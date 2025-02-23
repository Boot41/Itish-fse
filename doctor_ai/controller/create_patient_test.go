package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DoctorMock struct {
	ID    string `gorm:"type:text;primaryKey"`
	Email string `gorm:"unique"`
}

type PatientMock struct {
	ID   string `gorm:"type:text;primaryKey"`
	Name string
	Age  int
}

type TranscriptionMock struct {
	ID        string `gorm:"type:text;primaryKey"`
	DoctorID  string
	PatientID string
	Audio     string
	Report    string
}

var mockDBPatients *gorm.DB

// ✅ Setup mock DB for patients
func SetupMockDBPatients() {
	var err error
	mockDBPatients, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Run migrations
	mockDBPatients.AutoMigrate(&DoctorMock{}, &PatientMock{}, &TranscriptionMock{})

	// Insert mock doctor and patients
	mockDoctor := DoctorMock{ID: uuid.NewString(), Email: "doctor@example.com"}
	mockDBPatients.Create(&mockDoctor)

	mockPatient1 := PatientMock{ID: uuid.NewString(), Name: "John Doe", Age: 30}
	mockPatient2 := PatientMock{ID: uuid.NewString(), Name: "Jane Smith", Age: 25}
	mockDBPatients.Create(&mockPatient1)
	mockDBPatients.Create(&mockPatient2)

	mockTranscription := TranscriptionMock{
		ID:        uuid.NewString(),
		DoctorID:  mockDoctor.ID,
		PatientID: mockPatient1.ID,
		Audio:     "https://example.com/audio.mp3",
		Report:    "Medical report for John Doe",
	}
	mockDBPatients.Create(&mockTranscription)
}

// ✅ Mock handler for fetching patients
func MockGetPatients(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing email parameter"})
		return
	}

	// Handle DB failure
	if mockDBPatients == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database failure"})
		return
	}

	var doctor DoctorMock
	if err := mockDBPatients.Where("email = ?", email).First(&doctor).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Doctor not found"})
		return
	}

	var patients []PatientMock
	mockDBPatients.Find(&patients)

	c.JSON(http.StatusOK, patients)
}

// ✅ Test: Successfully fetch patients
func TestMockGetPatients_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBPatients()

	router := gin.Default()
	router.GET("/mock_get_patients", MockGetPatients)

	req, _ := http.NewRequest("GET", "/mock_get_patients?email=doctor@example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	assert.Contains(t, w.Body.String(), "Jane Smith")
}

// ✅ Test: Unauthorized Doctor (Doctor Not Found)
func TestMockGetPatients_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBPatients()

	router := gin.Default()
	router.GET("/mock_get_patients", MockGetPatients)

	req, _ := http.NewRequest("GET", "/mock_get_patients?email=unknown@example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized - Doctor not found")
}

// ✅ Test: Missing Email Parameter
func TestMockGetPatients_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBPatients()

	router := gin.Default()
	router.GET("/mock_get_patients", MockGetPatients)

	req, _ := http.NewRequest("GET", "/mock_get_patients", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing email parameter")
}

// ✅ Test: Simulated Database Failure
func TestMockGetPatients_DBFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBPatients()

	// Simulate DB failure
	mockDBPatients = nil

	router := gin.Default()
	router.GET("/mock_get_patients", MockGetPatients)

	req, _ := http.NewRequest("GET", "/mock_get_patients?email=doctor@example.com", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Database failure")
}
