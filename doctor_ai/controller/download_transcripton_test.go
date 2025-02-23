package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock database instance
var testDB *gorm.DB

// Doctor model
type TestDoctor struct {
	ID    string `gorm:"type:text;primaryKey"`
	Email string `gorm:"unique"`
}

// Patient model
type TestPatient struct {
	ID   string `gorm:"type:text;primaryKey"`
	Name string
}

// Transcription model
type TestTranscription struct {
	ID        string `gorm:"type:text;primaryKey"`
	DoctorID  string
	PatientID string
	Report    string
}

// Initialize mock SQLite DB
func setupTestDB() {
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Run migrations for test models
	testDB.AutoMigrate(&TestDoctor{}, &TestPatient{}, &TestTranscription{})

	// Insert test doctor
	testDB.Create(&TestDoctor{ID: "doc-123", Email: "doctor@example.com"})

	// Insert test patient
	testDB.Create(&TestPatient{ID: "pat-456", Name: "John Doe"})

	// Insert test transcription
	testDB.Create(&TestTranscription{
		ID:        "trans-789",
		DoctorID:  "doc-123",
		PatientID: "pat-456",
		Report:    "Test Medical Report",
	})
}

// Mock function for handling transcription downloads
func HandleDownloadTranscription(c *gin.Context) {
	var request struct {
		DoctorEmail string `json:"email"`
		PatientID   string `json:"patient_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Find doctor by email
	var doctor TestDoctor
	if err := testDB.Where("email = ?", request.DoctorEmail).First(&doctor).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Find transcription for the given doctor and patient
	var transcription TestTranscription
	if err := testDB.Where("doctor_id = ? AND patient_id = ?", doctor.ID, request.PatientID).First(&transcription).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transcription not found"})
		return
	}

	// Mock returning the transcription report as JSON
	c.JSON(http.StatusOK, gin.H{"report": transcription.Report})
}

// Unit test for downloading a transcription
func TestHandleDownloadTranscription_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB() // Set up test database

	router := gin.Default()
	router.POST("/download_transcription", HandleDownloadTranscription)

	requestBody := `{"email": "doctor@example.com", "patient_id": "pat-456"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/download_transcription", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)                     // Expect 200 OK
	assert.Contains(t, w.Body.String(), "Test Medical Report") // Verify report data
}
