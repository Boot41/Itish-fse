package service

import (
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTranscriptionTestDB(t *testing.T) (uuid.UUID, uuid.UUID, error) {
	// Use a unique database file for each test to avoid locking
	dbName := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create tables with SQLite-compatible types
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
		medical_history TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
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
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
		FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create transcriptions table: %v", err)
	}

	// Set the global DB instance
	initializers.DB = db

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
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		doctor_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id)
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create patients table: %v", err)
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS transcriptions (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		patient_id TEXT,
		text TEXT NOT NULL,
		report TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id),
		FOREIGN KEY (patient_id) REFERENCES patients(id)
	)`).Error
	if err != nil {
		t.Fatalf("Failed to create transcriptions table: %v", err)
	}

	// Create test doctor
	doctorID := uuid.New()
	doctor := &models.Doctor{
		ID:             doctorID,
		Name:          "Test Doctor",
		Email:         fmt.Sprintf("test%s@example.com", doctorID.String()[:8]),
		Password:      "hashedpassword",
		Phone:         fmt.Sprintf("1234%s", doctorID.String()[:6]),
		Specialization: "Test Specialization",
		CreatedAt:     time.Now(),
	}
	if err := db.Create(doctor).Error; err != nil {
		t.Fatalf("Failed to create test doctor: %v", err)
	}

	// Create test patient
	patientID := uuid.New()
	patient := &models.Patient{
		ID:        patientID,
		Name:     "Test Patient",
		Age:      30,
		Gender:   "Male",
		DoctorID: doctorID,
		CreatedAt: time.Now(),
	}
	if err := db.Create(patient).Error; err != nil {
		t.Fatalf("Failed to create test patient: %v", err)
	}

	return doctorID, patientID, nil
}

func createTestTranscription(t *testing.T, doctorID, patientID uuid.UUID) *models.Transcription {
	transcription := &models.Transcription{
		ID:        uuid.New(),
		DoctorID:  doctorID,
		PatientID: patientID,
		Text:      "Test transcription content",
		Report:    "Test transcription report",
		CreatedAt: time.Now(),
	}

	if err := initializers.DB.Create(transcription).Error; err != nil {
		t.Fatalf("Failed to create test transcription: %v", err)
	}

	return transcription
}

func TestGetTranscriptions(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "Valid request - first page",
			page:          1,
			limit:         3,
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name:          "Valid request - second page",
			page:          2,
			limit:         3,
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name:          "Invalid doctor ID",
			page:          1,
			limit:         10,
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fresh database for each test
			doctorID, patientID, err := setupTranscriptionTestDB(t)
			assert.NoError(t, err)

			// Create multiple test transcriptions
			for i := 0; i < 5; i++ {
				createTestTranscription(t, doctorID, patientID)
			}

			// Use a different doctor ID for the invalid case
			testDoctorID := doctorID
			if tt.name == "Invalid doctor ID" {
				testDoctorID = uuid.New()
			}

			// Get transcriptions
			transcriptions, _, total, err := GetTranscriptions(testDoctorID, tt.page, tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedCount > 0 {
					assert.Equal(t, int64(5), total) // Total should always be 5
					assert.Len(t, transcriptions, tt.expectedCount)
				} else {
					assert.Equal(t, int64(0), total)
					assert.Empty(t, transcriptions)
				}
			}
		})
	}
}

func TestGetTranscriptionByID(t *testing.T) {
	tests := []struct {
		name            string
		wantErr         bool
	}{
		{
			name:            "Valid transcription ID",
			wantErr:         false,
		},
		{
			name:            "Invalid transcription ID",
			wantErr:         true,
		},
	}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fresh database for each test
			doctorID, patientID, err := setupTranscriptionTestDB(t)
			assert.NoError(t, err)

			// Create test transcription
			transcription := createTestTranscription(t, doctorID, patientID)

			// Get transcription
			transcriptionID := transcription.ID
			if tt.wantErr {
				transcriptionID = uuid.New()
			}
			result, err := GetTranscriptionByID(transcriptionID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, transcription.ID, result.ID)
			}
		})
	}
}

func TestUpdateTranscription(t *testing.T) {
	doctorID, patientID, err := setupTranscriptionTestDB(t)
	assert.NoError(t, err)

	transcription := createTestTranscription(t, doctorID, patientID)

	tests := []struct {
		name            string
		transcriptionID uuid.UUID
		updateData      map[string]interface{}
		wantErr         bool
	}{
		{
			name:            "Valid update",
			transcriptionID: transcription.ID,
			updateData: map[string]interface{}{
				"Text": "Updated test content",
			},
			wantErr: false,
		},
		{
			name:            "Non-existent transcription ID",
			transcriptionID: uuid.New(),
			updateData: map[string]interface{}{
				"Text": "Updated test content",
			},
			wantErr: false, // Current service behavior: no error when updating non-existent record
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateTranscription(tt.transcriptionID, tt.updateData)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// For non-existent transcription, verify it wasn't created
				updated, err := GetTranscriptionByID(tt.transcriptionID)
				if tt.name == "Non-existent transcription ID" {
					assert.Error(t, err)
					assert.Nil(t, updated)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.updateData["Text"], updated.Text)
				}
			}
		})
	}
}

func TestDeleteTranscription(t *testing.T) {
	doctorID, patientID, err := setupTranscriptionTestDB(t)
	assert.NoError(t, err)

	transcription := createTestTranscription(t, doctorID, patientID)

	tests := []struct {
		name            string
		transcriptionID uuid.UUID
		wantErr         bool
	}{
		{
			name:            "Valid delete",
			transcriptionID: transcription.ID,
			wantErr:         false,
		},
		{
			name:            "Invalid transcription ID",
			transcriptionID: uuid.New(),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteTranscription(tt.transcriptionID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the deletion
				_, err := GetTranscriptionByID(tt.transcriptionID)
				assert.Error(t, err)
			}
		})
	}
}

func TestGetTranscriptionByPatient(t *testing.T) {
	doctorID, patientID, err := setupTranscriptionTestDB(t)
	assert.NoError(t, err)

	// Create multiple transcriptions for the same patient
	transcription1 := createTestTranscription(t, doctorID, patientID)
	time.Sleep(time.Millisecond) // Ensure different creation times
	transcription2 := createTestTranscription(t, doctorID, patientID)

	tests := []struct {
		name       string
		doctorID   uuid.UUID
		patientID  uuid.UUID
		wantErr    bool
		checkLatest bool
	}{
		{
			name:       "Valid request",
			doctorID:   doctorID,
			patientID:  patientID,
			wantErr:    false,
			checkLatest: true,
		},
		{
			name:       "Invalid doctor ID",
			doctorID:   uuid.New(),
			patientID:  patientID,
			wantErr:    true,
			checkLatest: false,
		},
		{
			name:       "Invalid patient ID",
			doctorID:   doctorID,
			patientID:  uuid.New(),
			wantErr:    true,
			checkLatest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTranscriptionByPatient(tt.doctorID, tt.patientID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if tt.checkLatest {
					// Should return the latest transcription (transcription2)
					assert.Equal(t, transcription2.ID, got.ID)
					assert.True(t, got.CreatedAt.After(transcription1.CreatedAt))
				}
			}
		})
	}
}
