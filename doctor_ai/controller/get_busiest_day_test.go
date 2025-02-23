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

func TestGetBusiestDays_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock busiest days data (simulating service layer)
	mockStats := []service.TimeBasedStats{
		{Date: "2025-01-25", Count: 100},
		{Date: "2025-01-24", Count: 90},
		{Date: "2025-01-23", Count: 80},
		{Date: "2025-01-22", Count: 70},
		{Date: "2025-01-21", Count: 60},
	}

	// Mock service function to return mock busiest days
	patches := gomonkey.ApplyFunc(service.GetBusiestDays, func(doctorID string, days int) ([]service.TimeBasedStats, error) {
		return mockStats, nil
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s"}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_busiest_days", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_busiest_days", GetBusiestDays)
	router.ServeHTTP(w, req)

	// Assertions: Check for OK status and the returned statistics
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2025-01-25")
	assert.Contains(t, w.Body.String(), "100")
	assert.Contains(t, w.Body.String(), "2025-01-24")
	assert.Contains(t, w.Body.String(), "90")
}

func TestGetBusiestDays_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_busiest_days", GetBusiestDays)

	// Missing email in the request body
	requestBody := `{"email": "", "days": 30}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_busiest_days", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for BadRequest status and error message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Doctor email is required")
}

func TestGetBusiestDays_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_busiest_days", GetBusiestDays)

	// Provide a non-existing doctor's email
	requestBody, _ := json.Marshal(map[string]interface{}{
		"email": "unknown@example.com",
		"days":  30,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_busiest_days", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for Unauthorized status and error message
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetBusiestDays_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock service function to simulate a failure
	patches := gomonkey.ApplyFunc(service.GetBusiestDays, func(doctorID string, days int) ([]service.TimeBasedStats, error) {
		return nil, errors.New("database error")
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 30}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_busiest_days", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_busiest_days", GetBusiestDays)
	router.ServeHTTP(w, req)

	// Assertions: Check for Internal Server Error and the appropriate error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve busiest days")
}

func TestGetBusiestDays_InvalidDaysQueryParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock busiest days data (simulating service layer)
	mockStats := []service.TimeBasedStats{
		{Date: "2025-01-25", Count: 100},
		{Date: "2025-01-24", Count: 90},
	}

	// Mock service function to return mock busiest days
	patches := gomonkey.ApplyFunc(service.GetBusiestDays, func(doctorID string, days int) ([]service.TimeBasedStats, error) {
		return mockStats, nil
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s"}`, uniqueEmail)

	// Invalid days query parameter
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_busiest_days?days=not_a_number", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_busiest_days", GetBusiestDays)
	router.ServeHTTP(w, req)

	// Assertions: Check that the default value is applied
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "2025-01-25")
	assert.Contains(t, w.Body.String(), "100")
}
