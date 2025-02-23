package service

import (
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPatientTestDB(t *testing.T) *gorm.DB {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clear existing data
	db.Exec("DELETE FROM doctors")
	db.Exec("DELETE FROM patients")
	db.Exec("DELETE FROM transcriptions")

	// Create SQLite-compatible tables
	err = db.Exec(`CREATE TABLE IF NOT EXISTS doctors (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		specialization TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		phone TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create doctors table: %v", err)
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS patients (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id)
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create patients table: %v", err)
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS transcriptions (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		patient_id TEXT NOT NULL,
		text TEXT NOT NULL,
		report TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id),
		FOREIGN KEY (patient_id) REFERENCES patients(id)
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create transcriptions table: %v", err)
	}

	// Set the global DB instance
	initializers.DB = db

	return db
}

func createTestDoctor(t *testing.T) *models.Doctor {
	doctor := &models.Doctor{
		ID:             uuid.New(),
		Name:           "Test Doctor",
		Email:          "test@example.com",
		Password:       "password123",
		Phone:          "1234567890",
		Specialization: "Test Specialization",
		CreatedAt:      time.Now(),
	}

	err := initializers.DB.Create(doctor).Error
	assert.NoError(t, err)
	return doctor
}

func TestValidatePatient(t *testing.T) {
	tests := []struct {
		name    string
		patient *models.Patient
		wantErr string
	}{
		{
			name: "Valid patient",
			patient: &models.Patient{
				Name:   "Test Patient",
				Age:    30,
				Gender: "Male",
			},
			wantErr: "",
		},
		{
			name: "Empty name",
			patient: &models.Patient{
				Age:    30,
				Gender: "Male",
			},
			wantErr: "patient name is required",
		},
		{
			name: "Invalid age",
			patient: &models.Patient{
				Name:   "Test Patient",
				Age:    0,
				Gender: "Male",
			},
			wantErr: "patient age must be greater than zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePatient(tt.patient)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPatients(t *testing.T) {
	db := setupPatientTestDB(t)
	defer db.Exec("DELETE FROM doctors")
	defer db.Exec("DELETE FROM patients")

	doctor := createTestDoctor(t)

	// Create test patients
	patients := []models.Patient{
		{
			ID:        uuid.New(),
			DoctorID:  doctor.ID,
			Name:      "Patient 1",
			Age:       30,
			Gender:    "Male",
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			DoctorID:  doctor.ID,
			Name:      "Patient 2",
			Age:       25,
			Gender:    "Female",
			CreatedAt: time.Now(),
		},
	}

	for _, p := range patients {
		err := initializers.DB.Create(&p).Error
		assert.NoError(t, err)
	}

	tests := []struct {
		name      string
		doctorID  uuid.UUID
		page      int
		limit     int
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "Valid pagination",
			doctorID:  doctor.ID,
			page:      1,
			limit:     1,
			wantCount: 1,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:      "Get all patients",
			doctorID:  doctor.ID,
			page:      1,
			limit:     10,
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:      "No patients for doctor",
			doctorID:  uuid.New(),
			page:      1,
			limit:     10,
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patients, total, err := GetPatients(tt.doctorID, tt.page, tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(patients))
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestGetPatientWithTranscription(t *testing.T) {
	db := setupPatientTestDB(t)
	defer db.Exec("DELETE FROM doctors")
	defer db.Exec("DELETE FROM patients")
	defer db.Exec("DELETE FROM transcriptions")

	doctor := createTestDoctor(t)

	// Create test patient
	patient := models.Patient{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		Name:      "Test Patient",
		Age:       30,
		Gender:    "Male",
		CreatedAt: time.Now(),
	}
	err := initializers.DB.Create(&patient).Error
	assert.NoError(t, err)

	// Create test transcription
	transcription := models.Transcription{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		PatientID: patient.ID,
		Text:      "Test transcription",
		Report:    "Test report",
		CreatedAt: time.Now(),
	}
	err = initializers.DB.Create(&transcription).Error
	assert.NoError(t, err)

	tests := []struct {
		name              string
		doctorID          uuid.UUID
		patientID         uuid.UUID
		wantPatient       bool
		wantTranscription bool
		wantErr           bool
	}{
		{
			name:              "Valid patient and transcription",
			doctorID:          doctor.ID,
			patientID:         patient.ID,
			wantPatient:       true,
			wantTranscription: true,
			wantErr:           false,
		},
		{
			name:              "Patient not found",
			doctorID:          doctor.ID,
			patientID:         uuid.New(),
			wantPatient:       false,
			wantTranscription: false,
			wantErr:           true,
		},
		{
			name:              "Unauthorized doctor",
			doctorID:          uuid.New(),
			patientID:         patient.ID,
			wantPatient:       false,
			wantTranscription: false,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patient, transcription, err := GetPatientWithTranscription(tt.doctorID, tt.patientID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, patient)
				assert.Nil(t, transcription)
			} else {
				assert.NoError(t, err)
				if tt.wantPatient {
					assert.NotNil(t, patient)
					assert.Equal(t, tt.patientID, patient.ID)
				}
				if tt.wantTranscription {
					assert.NotNil(t, transcription)
					assert.Equal(t, tt.patientID, transcription.PatientID)
				}
			}
		})
	}
}

func TestUpdatePatientByID(t *testing.T) {
	db := setupPatientTestDB(t)
	defer db.Exec("DELETE FROM doctors")
	defer db.Exec("DELETE FROM patients")

	doctor := createTestDoctor(t)

	// Create test patient
	patient := models.Patient{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		Name:      "Test Patient",
		Age:       30,
		Gender:    "Male",
		CreatedAt: time.Now(),
	}
	err := initializers.DB.Create(&patient).Error
	assert.NoError(t, err)

	tests := []struct {
		name       string
		doctorID   uuid.UUID
		patientID  uuid.UUID
		updateData map[string]interface{}
		wantErr    bool
	}{
		{
			name:      "Valid update",
			doctorID:  doctor.ID,
			patientID: patient.ID,
			updateData: map[string]interface{}{
				"name":   "Updated Name",
				"age":    35,
				"gender": "Female",
			},
			wantErr: false,
		},
		{
			name:      "Patient not found",
			doctorID:  doctor.ID,
			patientID: uuid.New(),
			updateData: map[string]interface{}{
				"name": "Updated Name",
			},
			wantErr: true,
		},
		{
			name:      "Unauthorized doctor",
			doctorID:  uuid.New(),
			patientID: patient.ID,
			updateData: map[string]interface{}{
				"name": "Updated Name",
			},
			wantErr: true,
		},
		{
			name:      "Disallowed fields",
			doctorID:  doctor.ID,
			patientID: patient.ID,
			updateData: map[string]interface{}{
				"name":      "Updated Name",
				"doctor_id": uuid.New(), // Should be ignored
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdatePatientByID(tt.doctorID, tt.patientID, tt.updateData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the update
				var updatedPatient models.Patient
				err := initializers.DB.First(&updatedPatient, "id = ?", tt.patientID).Error
				assert.NoError(t, err)

				// Check allowed fields were updated
				if name, ok := tt.updateData["name"].(string); ok {
					assert.Equal(t, name, updatedPatient.Name)
				}
				if age, ok := tt.updateData["age"].(int); ok {
					assert.Equal(t, age, updatedPatient.Age)
				}
				if gender, ok := tt.updateData["gender"].(string); ok {
					assert.Equal(t, gender, updatedPatient.Gender)
				}

				// Check disallowed fields were not updated
				if tt.name == "Disallowed fields" {
					assert.Equal(t, doctor.ID, updatedPatient.DoctorID)
				}
			}
		})
	}
}

func TestDeletePatient(t *testing.T) {
	db := setupPatientTestDB(t)
	defer db.Exec("DELETE FROM doctors")
	defer db.Exec("DELETE FROM patients")
	defer db.Exec("DELETE FROM transcriptions")

	doctor := createTestDoctor(t)

	// Create test patient
	patient := models.Patient{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		Name:      "Test Patient",
		Age:       30,
		Gender:    "Male",
		CreatedAt: time.Now(),
	}
	err := initializers.DB.Create(&patient).Error
	assert.NoError(t, err)

	// Create test transcription
	transcription := models.Transcription{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		PatientID: patient.ID,
		Text:      "Test transcription",
		Report:    "Test report",
		CreatedAt: time.Now(),
	}
	err = initializers.DB.Create(&transcription).Error
	assert.NoError(t, err)

	tests := []struct {
		name      string
		patientID uuid.UUID
		wantErr   bool
	}{
		{
			name:      "Valid deletion",
			patientID: patient.ID,
			wantErr:   false,
		},
		{
			name:      "Patient not found",
			patientID: uuid.New(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeletePatient(tt.patientID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify patient was deleted
				var deletedPatient models.Patient
				err := initializers.DB.First(&deletedPatient, "id = ?", tt.patientID).Error
				assert.Error(t, err)

				// Verify associated transcriptions were deleted
				var deletedTranscription models.Transcription
				err = initializers.DB.First(&deletedTranscription, "patient_id = ?", tt.patientID).Error
				assert.Error(t, err)
			}
		})
	}
}
