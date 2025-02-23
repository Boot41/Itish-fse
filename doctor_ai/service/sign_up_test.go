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

func setupSignUpTestDB(t *testing.T) *gorm.DB {
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

	// Set the global DB instance
	initializers.DB = db

	return db
}

func TestValidateDoctor(t *testing.T) {
	tests := []struct {
		name    string
		doctor  *models.Doctor
		wantErr string
	}{
		{
			name: "Valid doctor",
			doctor: &models.Doctor{
				Name:           "Test Doctor",
				Email:          "test@example.com",
				Password:       "password123",
				Phone:          "1234567890",
				Specialization: "Test Specialization",
			},
			wantErr: "",
		},
		{
			name: "Invalid email",
			doctor: &models.Doctor{
				Name:           "Test Doctor",
				Email:          "invalid-email",
				Password:       "password123",
				Phone:          "1234567890",
				Specialization: "Test Specialization",
			},
			wantErr: "invalid email format",
		},
		{
			name: "Short name",
			doctor: &models.Doctor{
				Name:           "T",
				Email:          "test@example.com",
				Password:       "password123",
				Phone:          "1234567890",
				Specialization: "Test Specialization",
			},
			wantErr: "name must be at least 2 characters long",
		},
		{
			name: "Invalid phone",
			doctor: &models.Doctor{
				Name:           "Test Doctor",
				Email:          "test@example.com",
				Password:       "password123",
				Phone:          "123",
				Specialization: "Test Specialization",
			},
			wantErr: "phone number must be 10 digits",
		},
		{
			name: "Short password",
			doctor: &models.Doctor{
				Name:           "Test Doctor",
				Email:          "test@example.com",
				Password:       "short",
				Phone:          "1234567890",
				Specialization: "Test Specialization",
			},
			wantErr: "password must be at least 8 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDoctor(tt.doctor)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateDoctor(t *testing.T) {
	tests := []struct {
		name    string
		doctor  *models.Doctor
		wantErr string
	}{
		{
			name: "Valid doctor creation",
			doctor: &models.Doctor{
				ID:             uuid.New(),
				Name:           "Test Doctor",
				Email:          "test@example.com",
				Password:       "password123",
				Phone:          "1234567890",
				Specialization: "Test Specialization",
				CreatedAt:      time.Now(),
			},
			wantErr: "",
		},
		{
			name: "Duplicate email",
			doctor: &models.Doctor{
				ID:             uuid.New(),
				Name:           "Another Doctor",
				Email:          "test@example.com", // Same email as previous test
				Password:       "password123",
				Phone:          "9876543210",
				Specialization: "Test Specialization",
				CreatedAt:      time.Now(),
			},
			wantErr: "email already exists",
		},
		{
			name: "Invalid input",
			doctor: &models.Doctor{
				ID:             uuid.New(),
				Name:           "T", // Invalid name
				Email:          "test2@example.com",
				Password:       "pass", // Invalid password
				Phone:          "123",  // Invalid phone
				Specialization: "Test Specialization",
				CreatedAt:      time.Now(),
			},
			wantErr: "name must be at least 2 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fresh database for each test
			_ = setupSignUpTestDB(t)

			// For duplicate email test, first create the initial doctor
			if tt.name == "Duplicate email" {
				initialDoctor := &models.Doctor{
					ID:             uuid.New(),
					Name:           "Test Doctor",
					Email:          "test@example.com",
					Password:       "password123",
					Phone:          "1234567890",
					Specialization: "Test Specialization",
					CreatedAt:      time.Now(),
				}
				err := CreateDoctor(initialDoctor)
				assert.NoError(t, err)
			}

			err := CreateDoctor(tt.doctor)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				// Verify the doctor was created
				var createdDoctor models.Doctor
				result := initializers.DB.Where("email = ?", tt.doctor.Email).First(&createdDoctor)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.doctor.Name, createdDoctor.Name)
				assert.Equal(t, tt.doctor.Email, createdDoctor.Email)
				assert.Equal(t, tt.doctor.Phone, createdDoctor.Phone)
				assert.Equal(t, tt.doctor.Specialization, createdDoctor.Specialization)

				// Verify the doctor was created with correct fields
				result = initializers.DB.Where("email = ?", tt.doctor.Email).First(&createdDoctor)
				assert.NoError(t, result.Error)

				// Check non-password fields
				assert.Equal(t, tt.doctor.Name, createdDoctor.Name)
				assert.Equal(t, tt.doctor.Email, createdDoctor.Email)
				assert.Equal(t, tt.doctor.Phone, createdDoctor.Phone)
				assert.Equal(t, tt.doctor.Specialization, createdDoctor.Specialization)

				// Verify password was hashed by checking it starts with bcrypt identifier
				assert.True(t, len(createdDoctor.Password) > 0)
				assert.Equal(t, "$2a$", createdDoctor.Password[:4], "Password should be hashed with bcrypt")
			}
		})
	}
}
