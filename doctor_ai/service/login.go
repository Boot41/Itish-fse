package service

import (
	"errors"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func DoctorLogin(user *models.Doctor) error {
	var existingUser models.Doctor

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

	// Find user by email
	result := tx.Where("email = ?", user.Email).First(&existingUser)

	if result.Error != nil {
		tx.Rollback()
		log.Println("Email not found...")
		return errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		tx.Rollback()
		log.Println("Incorrect password...")
		return errors.New("invalid email or password") // Always return generic error for security
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		return errors.New("login failed")
	}

	log.Println("Login successful:", existingUser.Email)
	return nil
}
