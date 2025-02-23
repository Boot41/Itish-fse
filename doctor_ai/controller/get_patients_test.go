package controller

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"itish41/doctor_ai_assistant/initializers"
)

type DoctorMockS struct {
	ID    string `gorm:"type:text;primaryKey"`
	Name  string
	Email string `gorm:"unique"`
}

type PatientMockS struct {
	ID       string `gorm:"type:text;primaryKey"`
	Name     string
	Age      int
	DoctorID string // ✅ Foreign key reference
}

type TranscriptionMockS struct {
	ID        string `gorm:"type:text;primaryKey"`
	DoctorID  string `gorm:"index"` // ✅ Foreign key index
	PatientID string `gorm:"index"` // ✅ Foreign key index
	Audio     string
	Report    string
}

func SetupMockDBForTests() {
	var err error
	initializers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to SQLite in-memory database")
	}

	// ✅ Ensure tables are created
	err = initializers.DB.AutoMigrate(&DoctorMockS{}, &PatientMockS{}, &TranscriptionMockS{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// ✅ Reset tables before inserting mock data
	initializers.DB.Exec("DELETE FROM patient_mock_s")
	initializers.DB.Exec("DELETE FROM doctor_mock_s")
	initializers.DB.Exec("DELETE FROM transcription_mock_s")

	// ✅ Insert a mock doctor
	mockDoctor := DoctorMockS{
		ID:    uuid.NewString(),
		Name:  "Dr. Mock",
		Email: "testdoctor@example.com",
	}
	initializers.DB.Create(&mockDoctor)

	log.Printf("Inserted Mock Doctor: %v", mockDoctor)

	// ✅ Insert mock patients linked to the doctor
	for i := 1; i <= 5; i++ {
		mockPatient := PatientMockS{
			ID:       uuid.NewString(),
			Name:     "Patient " + strconv.Itoa(i),
			Age:      i * 10,
			DoctorID: mockDoctor.ID, // ✅ Now properly linked
		}
		initializers.DB.Create(&mockPatient)
		log.Printf("Inserted Mock Patient: %v", mockPatient)
	}
}

// // ✅ Test case: Fetching patients successfully
// func TestGetPatients_Success(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	// Set up mock DB
// 	SetupMockDBForTests()

// 	// Create test router
// 	router := gin.Default()
// 	router.POST("/get_patients", GetPatients)

// 	// Test request body
// 	requestBody := `{"email": "testdoctor@example.com"}`
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/get_patients", strings.NewReader(requestBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)

// 	// ✅ Assert HTTP status is 200 OK
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// ✅ Assert response contains patients data
// 	assert.Contains(t, w.Body.String(), `"patients"`)
// }

// ✅ Test case: Doctor not found
func TestGetPatients_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up mock DB (empty this time)
	SetupMockDBForTests()

	// Create test router
	router := gin.Default()
	router.POST("/get_patients", GetPatients)

	// Test request with a non-existent email
	requestBody := `{"email": "nonexistent@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_patients", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// ✅ Assert HTTP status is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// ✅ Assert response message
	assert.Contains(t, w.Body.String(), `"Doctor not found"`)
}

// ✅ Test case: Missing email in request
func TestGetPatients_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up mock DB
	SetupMockDBForTests()

	// Create test router
	router := gin.Default()
	router.POST("/get_patients", GetPatients)

	// Test request with missing email
	requestBody := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_patients", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// ✅ Assert HTTP status is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// ✅ Assert response message
	assert.Contains(t, w.Body.String(), `"Email is required"`)
}
