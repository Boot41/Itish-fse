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

func setupProfileTestDB(t *testing.T) *gorm.DB {
	// Create a unique SQLite database for this test
	dbName := "file::memory:?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clear existing data
	db.Exec("DELETE FROM doctors")

	// Create SQLite-compatible table
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

	// Set the global DB instance
	initializers.DB = db

	return db
}

func createTestDoctorForProfile(t *testing.T, email string) *models.Doctor {
	doctor := &models.Doctor{
		ID:             uuid.New(),
		Name:           "Test Doctor",
		Email:          email,
		Password:       "password123",
		Phone:          "1234567890",
		Specialization: "Test Specialization",
		CreatedAt:      time.Now(),
	}

	err := initializers.DB.Create(doctor).Error
	assert.NoError(t, err)
	return doctor
}

func TestGetDoctorProfile(t *testing.T) {
	db := setupProfileTestDB(t)
	defer db.Exec("DELETE FROM doctors")

	doctor := createTestDoctorForProfile(t, "test@example.com")

	tests := []struct {
		name     string
		email    string
		wantName string
		wantErr  bool
	}{
		{
			name:     "Valid doctor",
			email:    doctor.Email,
			wantName: doctor.Name,
			wantErr:  false,
		},
		{
			name:     "Non-existent doctor",
			email:    "nonexistent@example.com",
			wantName: "",
			wantErr:  true,
		},
		{
			name:     "Empty email",
			email:    "",
			wantName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := GetDoctorProfile(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, tt.wantName, profile.Name)
				assert.Equal(t, tt.email, profile.Email)
				assert.Equal(t, doctor.Phone, profile.Phone)
				assert.Equal(t, doctor.Specialization, profile.Specialization)
				// Password should be cleared
				assert.Empty(t, profile.Password)
			}
		})
	}
}

func TestUpdateDoctorProfile(t *testing.T) {
	testEmail := "test@example.com" // Use a constant email for test cases

	tests := []struct {
		name              string
		email             string
		newName           string
		newPhone          string
		newSpecialization string
		wantErr           bool
	}{
		{
			name:              "Valid update - all fields",
			email:             testEmail,
			newName:           "Updated Doctor",
			newPhone:          "9876543210",
			newSpecialization: "Updated Specialization",
			wantErr:           false,
		},
		{
			name:              "Valid update - partial fields",
			email:             testEmail,
			newName:           "Another Name",
			newPhone:          "",
			newSpecialization: "",
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fresh database and doctor for each test case
			_ = setupProfileTestDB(t)
			testDoctor := createTestDoctorForProfile(t, tt.email)

			err := UpdateDoctorProfile(tt.email, tt.newName, tt.newSpecialization, tt.newPhone)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the update
				updatedProfile, err := GetDoctorProfile(tt.email)
				assert.NoError(t, err)
				assert.NotNil(t, updatedProfile)

				// Check if fields were updated correctly
				if tt.newName != "" {
					assert.Equal(t, tt.newName, updatedProfile.Name)
				}
				if tt.newPhone != "" {
					assert.Equal(t, tt.newPhone, updatedProfile.Phone)
				}
				if tt.newSpecialization != "" {
					assert.Equal(t, tt.newSpecialization, updatedProfile.Specialization)
				}

				// Fields that weren't updated should remain unchanged
				if tt.newName == "" {
					assert.Equal(t, testDoctor.Name, updatedProfile.Name)
				}
				if tt.newPhone == "" {
					assert.Equal(t, testDoctor.Phone, updatedProfile.Phone)
				}
				if tt.newSpecialization == "" {
					assert.Equal(t, testDoctor.Specialization, updatedProfile.Specialization)
				}
			}
		})
	}
}
