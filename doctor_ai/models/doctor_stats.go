package models

import (
	"time"

	"github.com/google/uuid"
)

type DoctorStats struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DoctorID      uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;"`
	Date          time.Time `gorm:"type:date;not null"`
	PatientsCount int       `gorm:"not null"`

	// Relationships
	Doctor Doctor `gorm:"foreignKey:DoctorID"`
	// gorm.Model
}
