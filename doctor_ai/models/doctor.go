package models

import (
	"time"

	"github.com/google/uuid"
)

type Doctor struct {
	// gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Specialization string    `gorm:"type:varchar(255);not null"`
	Email          string    `gorm:"type:varchar(255);unique;not null"`
	Password       string    `gorm:"type:varchar(255);not null"`
	Phone          string    `gorm:"type:varchar(20);unique;not null"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Patients       []Patient       `gorm:"foreignKey:DoctorID"`
	Transcriptions []Transcription `gorm:"foreignKey:DoctorID"`
	Statistics     []DoctorStats   `gorm:"foreignKey:DoctorID"`
}
