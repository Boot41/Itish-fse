package models

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Age       int       `gorm:"not null"`
	Gender    string    `gorm:"type:varchar(10);not null"`
	DoctorID  uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Doctor         Doctor          `gorm:"foreignKey:DoctorID"`
	Transcriptions []Transcription `gorm:"foreignKey:PatientID"`
	// gorm.Model
}
