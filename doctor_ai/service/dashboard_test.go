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

func setupTestDBSQLite(t *testing.T) (*uuid.UUID, *uuid.UUID) {
	var err error
	initializers.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create tables with SQLite-compatible schema
	err = initializers.DB.Exec(`
		CREATE TABLE IF NOT EXISTS doctors (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			specialization TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL DEFAULT '',
			phone TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS patients (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER NOT NULL,
			gender TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id)
		);

		CREATE TABLE IF NOT EXISTS transcriptions (
			id TEXT PRIMARY KEY,
			doctor_id TEXT NOT NULL,
			patient_id TEXT NOT NULL,
			text TEXT NOT NULL,
			report TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id),
			FOREIGN KEY (patient_id) REFERENCES patients(id)
		);
	`).Error
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Create test doctor
	doctorID := uuid.New()
	doctor := models.Doctor{
		ID:    doctorID,
		Name:  "Test Doctor",
		Email: "test@example.com",
	}
	if err := initializers.DB.Create(&doctor).Error; err != nil {
		t.Fatalf("Failed to create test doctor: %v", err)
	}

	// Create test patient
	patientID := uuid.New()
	patient := models.Patient{
		ID:       patientID,
		Name:     "Test Patient",
		Age:      30,
		Gender:   "Male",
		DoctorID: doctorID,
	}
	if err := initializers.DB.Create(&patient).Error; err != nil {
		t.Fatalf("Failed to create test patient: %v", err)
	}

	return &doctorID, &patientID
}

func createTestTranscriptions(t *testing.T, doctorID uuid.UUID, patientID uuid.UUID, count int) {
	// Insert transcriptions directly using SQL to ensure proper date handling
	for i := 0; i < count; i++ {
		// Calculate date for each transcription
		createdAt := time.Now().AddDate(0, 0, -i).UTC().Format("2006-01-02 15:04:05")

		// Insert using SQL with proper date format
		err := initializers.DB.Exec(`
			INSERT INTO transcriptions (id, doctor_id, patient_id, text, report, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			uuid.New().String(),
			doctorID.String(),
			patientID.String(),
			"Test transcription",
			"Test report",
			createdAt,
		).Error

		if err != nil {
			t.Fatalf("Failed to create test transcription: %v", err)
		}
	}
}

func TestGetDashboardTranscripts(t *testing.T) {
	t.Run("Valid request within range", func(t *testing.T) {
		doctorID, patientID := setupTestDBSQLite(t)
		createTestTranscriptions(t, *doctorID, *patientID, 5)

		transcriptions, err := getDashboardTranscriptsForTest(doctorID.String(), 7)
		assert.NoError(t, err)
		assert.NotEmpty(t, transcriptions)
		assert.LessOrEqual(t, len(transcriptions), 10)
	})

	t.Run("Invalid doctor ID", func(t *testing.T) {
		_, err := getDashboardTranscriptsForTest("invalid", 7)
		assert.Error(t, err)
	})

	t.Run("Zero days", func(t *testing.T) {
		doctorID, patientID := setupTestDBSQLite(t)
		createTestTranscriptions(t, *doctorID, *patientID, 5)

		transcriptions, err := getDashboardTranscriptsForTest(doctorID.String(), 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, transcriptions)
	})
}

func TestGetDashboardStats(t *testing.T) {
	t.Run("Valid doctor ID", func(t *testing.T) {
		doctorID, patientID := setupTestDBSQLite(t)
		createTestTranscriptions(t, *doctorID, *patientID, 5)

		stats, err := getDashboardStatsForTest(doctorID.String())
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(1), stats.TotalPatients)
		assert.Equal(t, int64(5), stats.TotalTranscriptions)
		assert.NotEmpty(t, stats.RecentTranscriptions)
		assert.NotEmpty(t, stats.RecentPatients)
	})
}

func TestGetDailyStatistics(t *testing.T) {
	doctorID, patientID := setupTestDBSQLite(t)
	createTestTranscriptions(t, *doctorID, *patientID, 5)

	tests := []struct {
		name        string
		doctorID    string
		days        int
		expectError bool
	}{
		{
			name:        "Valid request",
			doctorID:    doctorID.String(),
			days:        7,
			expectError: false,
		},
		{
			name:        "Invalid doctor ID",
			doctorID:    "invalid-uuid",
			days:        7,
			expectError: true,
		},
		{
			name:        "Zero days",
			doctorID:    doctorID.String(),
			days:        0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := getDailyStatisticsForTest(tt.doctorID, tt.days)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.days > 0 {
					assert.NotEmpty(t, stats)
				}
			}
		})
	}
}

func TestGetMonthlyStatistics(t *testing.T) {
	doctorID, patientID := setupTestDBSQLite(t)
	createTestTranscriptions(t, *doctorID, *patientID, 5)

	tests := []struct {
		name        string
		doctorID    string
		months      int
		expectError bool
	}{
		{
			name:        "Valid request",
			doctorID:    doctorID.String(),
			months:      3,
			expectError: false,
		},
		{
			name:        "Invalid doctor ID",
			doctorID:    "invalid-uuid",
			months:      3,
			expectError: true,
		},
		{
			name:        "Zero months",
			doctorID:    doctorID.String(),
			months:      0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := getMonthlyStatisticsForTest(tt.doctorID, tt.months)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.months > 0 {
					assert.NotEmpty(t, stats)
				}
			}
		})
	}
}

func TestGetBusiestDays(t *testing.T) {
	doctorID, patientID := setupTestDBSQLite(t)
	createTestTranscriptions(t, *doctorID, *patientID, 5)

	tests := []struct {
		name        string
		doctorID    string
		days        int
		expectError bool
	}{
		{
			name:        "Valid request",
			doctorID:    doctorID.String(),
			days:        7,
			expectError: false,
		},
		{
			name:        "Invalid doctor ID",
			doctorID:    "invalid-uuid",
			days:        7,
			expectError: true,
		},
		{
			name:        "Zero days",
			doctorID:    doctorID.String(),
			days:        0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := getBusiestDaysForTest(tt.doctorID, tt.days)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.days > 0 {
					assert.NotEmpty(t, stats)
				}
			}
		})
	}
}

func TestGetDashboardPatients(t *testing.T) {
	doctorID, _ := setupTestDBSQLite(t)

	// Create additional test patients using direct SQL
	for i := 0; i < 5; i++ {
		// Calculate date for each patient
		createdAt := time.Now().AddDate(0, 0, -i).UTC().Format("2006-01-02 15:04:05")

		// Insert using SQL with proper date format
		err := initializers.DB.Exec(`
			INSERT INTO patients (id, name, age, gender, doctor_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			uuid.New().String(),
			"Test Patient",
			30,
			"Male",
			doctorID.String(),
			createdAt,
		).Error

		if err != nil {
			t.Fatalf("Failed to create test patient: %v", err)
		}
	}

	tests := []struct {
		name          string
		doctorID      string
		days          int
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Valid request",
			doctorID:      doctorID.String(),
			days:          7,
			expectedCount: 6, // Including the one created in setupTestDB
			expectError:   false,
		},
		{
			name:          "Invalid doctor ID",
			doctorID:      "invalid-uuid",
			days:          7,
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "Zero days",
			doctorID:      doctorID.String(),
			days:          0,
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patients, err := getDashboardPatientsForTest(tt.doctorID, tt.days)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, patients, tt.expectedCount)
			}
		})
	}
}
