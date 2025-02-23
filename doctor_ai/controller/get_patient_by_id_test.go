package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
)

var mockDB *gorm.DB

func SetupMockDBForPatientTest() {
	var err error
	mockDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to SQLite in-memory database")
	}

	mockDB.Exec(`
	CREATE TABLE doctors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		specialization TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		phone TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	mockDB.Exec(`
	CREATE TABLE patients (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		doctor_id TEXT NOT NULL,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id)
	);
	`)
	mockDB.Exec(`
	CREATE TABLE transcriptions (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		patient_id TEXT NOT NULL,
		text TEXT NOT NULL,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id),
		FOREIGN KEY (patient_id) REFERENCES patients(id)
	);
	`)

	initializers.DB = mockDB

	doctorID := uuid.NewString()
	patientID := uuid.NewString()
	transcriptionID := uuid.NewString()

	mockDB.Exec(`INSERT INTO doctors (id, name, specialization, email, password, phone) VALUES (?, ?, ?, ?, ?, ?)`,
		doctorID, "Test Doctor", "Cardiology", "testdoctor@example.com", "password123", "1234567890")

	mockDB.Exec(`INSERT INTO patients (id, name, doctor_id) VALUES (?, ?, ?)`,
		patientID, "Test Patient", doctorID)

	mockDB.Exec(`INSERT INTO transcriptions (id, doctor_id, patient_id, text) VALUES (?, ?, ?, ?)`,
		transcriptionID, doctorID, patientID, "Test transcription")

	var tables []string
	mockDB.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
	fmt.Println("Tables in DB:", tables)
}

func TestGetPatientByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBForPatientTest()

	router := gin.Default()
	router.POST("/get_patient", GetPatientByID)

	var patient models.Patient
	mockDB.First(&patient)

	requestBody := `{"email": "testdoctor@example.com", "patient_id": "` + patient.ID.String() + `"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_patient", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"patient"`)
	assert.Contains(t, w.Body.String(), `"transcription"`)
}

// func TestGetPatientByID_InvalidPatientID(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	SetupMockDBForPatientTest()

// 	router := gin.Default()
// 	router.POST("/get_patient", GetPatientByID)

// 	invalidPatientID := uuid.NewString()
// 	requestBody := `{"email": "testdoctor@example.com", "patient_id": "` + invalidPatientID + `"}`

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/get_patient", strings.NewReader(requestBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	assert.Contains(t, w.Body.String(), `"message":"Patient not found"`)
// }

func TestGetPatientByID_InvalidDoctorEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupMockDBForPatientTest()

	router := gin.Default()
	router.POST("/get_patient", GetPatientByID)

	var patient models.Patient
	mockDB.First(&patient)

	requestBody := `{"email": "wrongdoctor@example.com", "patient_id": "` + patient.ID.String() + `"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_patient", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Doctor not found"`)
}
