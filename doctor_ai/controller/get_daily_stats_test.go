package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetDailyStatistics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock daily statistics data (simulating service layer)
	mockStats := []service.TimeBasedStats{
		{Date: "2025-02-01", Count: 5},
		{Date: "2025-02-02", Count: 3},
	}

	// Mock service function to return mock statistics
	patches := gomonkey.ApplyFunc(service.GetDailyStatistics, func(doctorID string, days int) ([]service.TimeBasedStats, error) {
		return mockStats, nil
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 7}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_daily_statistics", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_daily_statistics", GetDailyStatistics)
	router.ServeHTTP(w, req)

	// Assertions: Check for OK status and the returned statistics
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2025-02-01")
	assert.Contains(t, w.Body.String(), "5")
	assert.Contains(t, w.Body.String(), "2025-02-02")
	assert.Contains(t, w.Body.String(), "3")
}

func TestGetDailyStatistics_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_daily_statistics", GetDailyStatistics)

	// Missing email in the request body
	requestBody := `{"email": "", "days": 7}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_daily_statistics", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for BadRequest status and error message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Doctor email is required")
}

func TestGetDailyStatistics_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_daily_statistics", GetDailyStatistics)

	// Provide a non-existing doctor's email
	requestBody, _ := json.Marshal(map[string]interface{}{
		"email": "unknown@example.com",
		"days":  7,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_daily_statistics", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for Unauthorized status and error message
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetDailyStatistics_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock service function to simulate a failure
	patches := gomonkey.ApplyFunc(service.GetDailyStatistics, func(doctorID string, days int) ([]service.TimeBasedStats, error) {
		return nil, errors.New("database error")
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 7}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_daily_statistics", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_daily_statistics", GetDailyStatistics)
	router.ServeHTTP(w, req)

	// Assertions: Check for Internal Server Error and the appropriate error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve daily statistics")
}
