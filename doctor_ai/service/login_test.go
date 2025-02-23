package service

import (
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupLoginTestDB(t *testing.T) (*models.Doctor, error) {
	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set the global DB instance
	initializers.DB = db

	// Create SQLite-compatible table structure
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

	// Create a test doctor with a known password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	testDoctor := &models.Doctor{
		ID:             uuid.New(),
		Email:          "test@example.com",
		Password:       string(hashedPassword),
		Name:           "Test Doctor",
		Phone:          "1234567890",
		Specialization: "Test Specialization",
		CreatedAt:      time.Now(),
	}

	if err := db.Create(testDoctor).Error; err != nil {
		t.Fatalf("Failed to create test doctor: %v", err)
	}

	return testDoctor, nil
}

func TestDoctorLogin(t *testing.T) {
	testDoctor, err := setupLoginTestDB(t)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		user    *models.Doctor
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid credentials",
			user: &models.Doctor{
				Email:    testDoctor.Email,
				Password: "testpassword",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			user: &models.Doctor{
				Email:    "nonexistent@example.com",
				Password: "testpassword",
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "Invalid password",
			user: &models.Doctor{
				Email:    testDoctor.Email,
				Password: "wrongpassword",
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "Empty credentials",
			user: &models.Doctor{
				Email:    "",
				Password: "",
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DoctorLogin(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
