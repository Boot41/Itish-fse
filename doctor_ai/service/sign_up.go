package service

import (
	"errors"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func validateDoctor(doctor *models.Doctor) error {
	// Email validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(doctor.Email) {
		return errors.New("invalid email format")
	}

	// Name validation
	if len(doctor.Name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	// Phone validation (basic check)
	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(doctor.Phone) {
		return errors.New("phone number must be 10 digits")
	}

	// Password validation
	if len(doctor.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func CreateDoctor(doctor *models.Doctor) error {
	// Validate doctor input
	if err := validateDoctor(doctor); err != nil {
		log.Println("Validation error:", err)
		return err
	}

	// Create a clean session with prepared statements disabled
	db := initializers.DB.Session(&gorm.Session{
		PrepareStmt:            false,
		SkipDefaultTransaction: true,
	})

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if email already exists
	var existingDoctor models.Doctor
	result := tx.Where("email = ?", doctor.Email).First(&existingDoctor)
	if result.Error == nil {
		tx.Rollback()
		return errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		log.Println("Password hashing error:", err)
		return errors.New("error processing password")
	}

	// Replace plain password with hashed password
	doctor.Password = string(hashedPassword)

	// Create the doctor in the database
	result = tx.Create(doctor)
	if result.Error != nil {
		tx.Rollback()
		log.Println("Database creation error:", result.Error)
		return errors.New("failed to create doctor")
	}

	tx.Commit()
	log.Println("Doctor created successfully:", doctor.Email)
	return nil
}
