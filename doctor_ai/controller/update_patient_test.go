package controller

import (
	"bytes"
	"encoding/json"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/service"
	"net/http/httptest"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func settingupTestDB(t *testing.T) {
	var err error
	initializers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create doctors table
	initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS doctors (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		);
	`)

	// Create patients table
	initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS patients (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id)
		);
	`)
}

func TestUpdatePatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	settingupTestDB(t)

	// Insert test doctor
	doctorID := uuid.New()
	initializers.DB.Exec(`INSERT INTO doctors (id, name, email) VALUES (?, ?, ?)`,
		doctorID.String(), "Test Doctor", "johndoe@example.com")

	// Insert test patient
	patientID := uuid.New()
	initializers.DB.Exec(`INSERT INTO patients (id, name, age, gender, doctor_id) VALUES (?, ?, ?, ?, ?)`,
		patientID.String(), "Old Name", 25, "Male", doctorID.String())

	// Create a test patient update request
	requestData := map[string]interface{}{
		"email":      "johndoe@example.com",
		"patient_id": patientID.String(),
		"update_data": map[string]interface{}{
			"name":   "John Doe",
			"age":    30,
			"gender": "Male",
		},
	}
	reqBody, _ := json.Marshal(requestData)

	// Create a POST request to update patient
	req := httptest.NewRequest("POST", "/update_patient", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Set up router and register the endpoint
	r := gin.Default()
	r.POST("/update_patient", UpdatePatient)

	// Mock the service layer UpdatePatientByID function
	patches := gomonkey.ApplyFunc(service.UpdatePatientByID, func(doctorID uuid.UUID, patientID uuid.UUID, updateData map[string]interface{}) error {
		return nil
	})
	defer patches.Reset()

	// Serve the request
	r.ServeHTTP(w, req)

	// Assert status code and response
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Patient updated successfully")
}
