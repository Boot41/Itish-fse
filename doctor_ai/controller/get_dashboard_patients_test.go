package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
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

func TestGetDashboardPatients_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	mockDoctorUUID := uuid.Must(uuid.Parse(mockDoctorID)) // Convert string to uuid.UUID
	// Insert mock patients (simulating service layer)
	mockPatients := []models.Patient{
		{ID: uuid.New(), Name: "Patient 1", DoctorID: mockDoctorUUID},
		{ID: uuid.New(), Name: "Patient 2", DoctorID: mockDoctorUUID},
	}
	// Mock service function to return mock patients
	patches := gomonkey.ApplyFunc(service.GetDashboardPatients, func(doctorID string, days int) ([]models.Patient, error) {
		return mockPatients, nil
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 7}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_patients", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_dashboard_patients", GetDashboardPatients)
	router.ServeHTTP(w, req)

	// Assertions: Check for OK status and the returned patients
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Patient 1")
	assert.Contains(t, w.Body.String(), "Patient 2")
}

func TestGetDashboardPatients_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_dashboard_patients", GetDashboardPatients)

	// Missing email in the request body
	requestBody := `{"email": "", "days": 7}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_patients", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for BadRequest status and error message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Doctor email is required")
}

func TestGetDashboardPatients_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_dashboard_patients", GetDashboardPatients)

	// Provide a non-existing doctor's email
	requestBody, _ := json.Marshal(map[string]interface{}{
		"email": "unknown@example.com",
		"days":  7,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_patients", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions: Check for Unauthorized status and error message
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetDashboardPatients_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock service function to simulate a failure
	patches := gomonkey.ApplyFunc(service.GetDashboardPatients, func(doctorID string, days int) ([]models.Patient, error) {
		return nil, errors.New("database error")
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 7}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_patients", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_dashboard_patients", GetDashboardPatients)
	router.ServeHTTP(w, req)

	// Assertions: Check for Internal Server Error and the appropriate error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve dashboard patients")
}
