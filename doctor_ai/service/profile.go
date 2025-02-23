package service

import (
	"database/sql"
	"errors"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"

	"gorm.io/gorm"
)

func GetDoctorProfile(email string) (*models.Doctor, error) {
	var doctor models.Doctor

	// Use database/sql directly to bypass GORM's prepared statements
	sqlDB, err := initializers.DB.DB()
	if err != nil {
		log.Println("Error getting database connection:", err)
		return nil, errors.New("database connection error")
	}

	// Directly query the database
	query := "SELECT id, email, name, phone, specialization, password FROM doctors WHERE email = $1 ORDER BY id LIMIT 1"
	err = sqlDB.QueryRow(query, email).Scan(
		&doctor.ID,
		&doctor.Email,
		&doctor.Name,
		&doctor.Phone,
		&doctor.Specialization,
		&doctor.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No doctor found with email:", email)
			return nil, errors.New("doctor profile not found")
		}
		log.Println("Error retrieving doctor profile:", err)
		return nil, errors.New("error retrieving doctor profile")
	}

	log.Println("Doctor profile retrieved successfully:", doctor.Email)

	// Clear sensitive information
	doctor.Password = ""
	return &doctor, nil
}

func UpdateDoctorProfile(email, name, specialization, phone string) error {
	var doctor models.Doctor

	// Create a new database session
	db := initializers.DB.Session(&gorm.Session{PrepareStmt: false})

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find doctor by email
	result := tx.Where("email = ?", email).First(&doctor)
	if result.Error != nil {
		log.Println("Doctor not found:", result.Error)
		return errors.New("doctor profile not found")
	}

	// Prepare fields to update
	updateData := map[string]interface{}{}
	if name != "" {
		updateData["Name"] = name
	}
	if specialization != "" {
		updateData["Specialization"] = specialization
	}
	if phone != "" {
		updateData["Phone"] = phone
	}

	// Update the doctor profile
	result = tx.Model(&doctor).Updates(updateData)
	if result.Error != nil {
		tx.Rollback()
		log.Println("Error updating doctor profile:", result.Error)
		return errors.New("failed to update doctor profile")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		return errors.New("failed to update doctor profile")
	}

	log.Println("Doctor profile updated successfully:", doctor.Email)
	return nil
}
