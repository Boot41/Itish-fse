package models

import (
	"time"

	"github.com/google/uuid"
)

type Transcription struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DoctorID  uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;"`
	PatientID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;"`
	Text      string    `gorm:"type:text;not null"`
	Report    string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Doctor  Doctor  `gorm:"foreignKey:DoctorID"`
	Patient Patient `gorm:"foreignKey:PatientID"`
	// gorm.Model
}
