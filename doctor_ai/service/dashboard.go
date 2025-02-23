package service

import (
	"log"

	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DashboardStats represents summary statistics for a doctor.
type DashboardStats struct {
	TotalPatients        int64
	TotalTranscriptions  int64
	RecentTranscriptions []models.Transcription
	RecentPatients       []models.Patient
}

func GetDashboardTranscripts(doctorID string, days int) ([]models.Transcription, error) {
	var transcriptions []models.Transcription

	// Convert string to UUID
	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Where("doctor_id = ? AND created_at >= CURRENT_DATE - ? * INTERVAL '1 day'", doctorUUID, days).
		Order("created_at DESC").
		Limit(10).
		Find(&transcriptions)
	if result.Error != nil {
		log.Println("Error retrieving dashboard transcripts:", result.Error)
		return nil, result.Error
	}

	return transcriptions, nil
}

func GetDashboardStats(doctorID string) (*DashboardStats, error) {
	var stats DashboardStats

	// Convert doctorID to UUID
	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	// Total Patients
	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Patient{}).
		Where("doctor_id = ?", doctorUUID).
		Count(&stats.TotalPatients)
	if result.Error != nil {
		log.Println("Error counting patients:", result.Error)
		return nil, result.Error
	}

	// Total Transcriptions
	result = initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Where("doctor_id = ?", doctorUUID).
		Count(&stats.TotalTranscriptions)
	if result.Error != nil {
		log.Println("Error counting transcriptions:", result.Error)
		return nil, result.Error
	}

	// Recent Transcriptions
	trans, err := GetDashboardTranscripts(doctorID, 30)
	if err != nil {
		return nil, err
	}
	stats.RecentTranscriptions = trans

	// Recent Patients
	patients, err := GetDashboardPatients(doctorID, 30)
	if err != nil {
		return nil, err
	}
	stats.RecentPatients = patients

	return &stats, nil
}

// TimeBasedStats represents daily or monthly statistics.
type TimeBasedStats struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

func GetDailyStatistics(doctorID string, days int) ([]TimeBasedStats, error) {
	var dailyStats []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Select("to_char(DATE(created_at), 'YYYY-MM-DD') as date, CAST(COUNT(*) AS BIGINT) as count").
		Where("doctor_id = ? AND created_at >= CURRENT_DATE - ? * INTERVAL '1 day'", doctorUUID, days).
		Group("to_char(DATE(created_at), 'YYYY-MM-DD')").
		Order("date DESC").
		Scan(&dailyStats)
	if result.Error != nil {
		log.Println("Error retrieving daily statistics:", result.Error)
		return nil, result.Error
	}

	return dailyStats, nil
}

func GetMonthlyStatistics(doctorID string, months int) ([]TimeBasedStats, error) {
	var monthlyStats []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Select("to_char(created_at, 'YYYY-MM') as date, CAST(COUNT(*) AS BIGINT) as count").
		Where("doctor_id = ? AND created_at >= CURRENT_DATE - ? * INTERVAL '1 month'", doctorUUID, months).
		Group("to_char(created_at, 'YYYY-MM')").
		Order("date DESC").
		Scan(&monthlyStats)
	if result.Error != nil {
		log.Println("Error retrieving monthly statistics:", result.Error)
		return nil, result.Error
	}

	return monthlyStats, nil
}

func GetBusiestDays(doctorID string, days int) ([]TimeBasedStats, error) {
	var busiestDays []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Select("to_char(DATE(created_at), 'YYYY-MM-DD') as date, CAST(COUNT(*) AS BIGINT) as count").
		Where("doctor_id = ? AND created_at >= CURRENT_DATE - ? * INTERVAL '1 day'", doctorUUID, days).
		Group("to_char(DATE(created_at), 'YYYY-MM-DD')").
		Order("count DESC").
		Limit(5).
		Scan(&busiestDays)
	if result.Error != nil {
		log.Println("Error retrieving busiest days:", result.Error)
		return nil, result.Error
	}

	return busiestDays, nil
}

func GetDashboardPatients(doctorID string, days int) ([]models.Patient, error) {
	var patients []models.Patient

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	query := `
		SELECT id, name, age, gender, doctor_id, created_at
		FROM patients
		WHERE doctor_id = ? AND created_at >= CURRENT_DATE - ? * INTERVAL '1 day'
		ORDER BY created_at DESC
		LIMIT 10
	`

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).Raw(query, doctorUUID, days).Scan(&patients)
	if result.Error != nil {
		log.Println("Error retrieving dashboard patients:", result.Error)
		return nil, result.Error
	}

	return patients, nil
}
