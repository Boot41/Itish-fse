package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMockDBForDashboardTest() {
	var err error
	initializers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to SQLite in-memory database")
	}

	initializers.DB.Exec(`
	CREATE TABLE doctors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);
	`)

	// Insert mock doctor with unique email for each test
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)
}

func TestGetDashboardTranscripts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup mock DB
	setupMockDBForDashboardTest()

	// Mock service function to avoid actual database call
	patches := gomonkey.ApplyFunc(service.GetDashboardTranscripts, func(doctorID string, days int) ([]models.Transcription, error) {
		return []models.Transcription{
			{ID: uuid.New(), DoctorID: uuid.New(), PatientID: uuid.New(), Text: "Sample text", Report: "Sample report", CreatedAt: time.Now()},
		}, nil
	})
	defer patches.Reset()

	// Insert mock doctor with a unique email for the test
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Use unique email in the request body
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 30}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_transcripts", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_dashboard_transcripts", GetDashboardTranscripts)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDashboardTranscripts_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_dashboard_transcripts", GetDashboardTranscripts)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"email": "",
		"days":  7,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_transcripts", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Doctor email is required")
}

func TestGetDashboardTranscripts_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	router := gin.Default()
	router.POST("/get_dashboard_transcripts", GetDashboardTranscripts)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"email": "unknown@example.com",
		"days":  7,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_transcripts", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestGetDashboardTranscripts_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockDBForDashboardTest()

	// Insert mock doctor with a unique email
	mockDoctorID := uuid.NewString()
	uniqueEmail := fmt.Sprintf("testdoctor_%s@example.com", uuid.NewString())
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		mockDoctorID, "Test Doctor", uniqueEmail)

	// Mock service function to simulate a database error
	patches := gomonkey.ApplyFunc(service.GetDashboardTranscripts, func(doctorID string, days int) ([]models.Transcription, error) {
		return nil, errors.New("database error")
	})

	defer patches.Reset()

	// Create request body with the email of the inserted mock doctor
	requestBody := fmt.Sprintf(`{"email": "%s", "days": 7}`, uniqueEmail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_dashboard_transcripts", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router := gin.Default()
	router.POST("/get_dashboard_transcripts", GetDashboardTranscripts)
	router.ServeHTTP(w, req)

	// Assertions: Check for Internal Server Error and appropriate error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve dashboard transcripts")
}
