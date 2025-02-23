package service

import (
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SQLite-specific version of GetDashboardTranscripts for testing
func getDashboardTranscriptsForTest(doctorID string, days int) ([]models.Transcription, error) {
	var transcriptions []models.Transcription

	// Convert string to UUID
	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		log.Println("Error parsing doctor UUID:", err)
		return nil, err
	}

	// Calculate the date threshold using SQLite's datetime function
	dateThreshold := time.Now().AddDate(0, 0, -days).UTC().Format("2006-01-02 15:04:05")

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&models.Transcription{}).
		Where("doctor_id = ? AND created_at >= ?", doctorUUID, dateThreshold).
		Order("created_at DESC").
		Limit(10).
		Find(&transcriptions)
	if result.Error != nil {
		log.Println("Error retrieving dashboard transcripts:", result.Error)
		return nil, result.Error
	}

	return transcriptions, nil
}

// SQLite-specific version of GetDashboardStats for testing
func getDashboardStatsForTest(doctorID string) (*DashboardStats, error) {
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
	trans, err := getDashboardTranscriptsForTest(doctorID, 30)
	if err != nil {
		return nil, err
	}
	stats.RecentTranscriptions = trans

	// Recent Patients
	patients, err := getDashboardPatientsForTest(doctorID, 30)
	if err != nil {
		return nil, err
	}
	stats.RecentPatients = patients

	return &stats, nil
}

// SQLite-specific version of GetDailyStatistics for testing
func getDailyStatisticsForTest(doctorID string, days int) ([]TimeBasedStats, error) {
	var stats []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %v", err)
	}

	// Calculate the date threshold using SQLite's datetime function
	dateThreshold := time.Now().AddDate(0, 0, -days).UTC().Format("2006-01-02")

	// Use strftime for date formatting in SQLite
	result := initializers.DB.Raw(`
		SELECT strftime('%Y-%m-%d', created_at) as date, COUNT(*) as count
		FROM transcriptions
		WHERE doctor_id = ?
		AND date(created_at) >= date(?)
		GROUP BY date
		ORDER BY date DESC
	`, doctorUUID, dateThreshold).Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

// SQLite-specific version of GetMonthlyStatistics for testing
func getMonthlyStatisticsForTest(doctorID string, months int) ([]TimeBasedStats, error) {
	var stats []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %v", err)
	}

	// Calculate the date threshold using SQLite's datetime function
	dateThreshold := time.Now().AddDate(0, -months, 0).UTC().Format("2006-01-02")

	// Use strftime for month formatting in SQLite
	result := initializers.DB.Raw(`
		SELECT strftime('%Y-%m', created_at) as date, COUNT(*) as count
		FROM transcriptions
		WHERE doctor_id = ?
		AND date(created_at) >= date(?)
		GROUP BY strftime('%Y-%m', created_at)
		ORDER BY date DESC
	`, doctorUUID, dateThreshold).Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

// SQLite-specific version of GetBusiestDays for testing
func getBusiestDaysForTest(doctorID string, days int) ([]TimeBasedStats, error) {
	var stats []TimeBasedStats

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %v", err)
	}

	// Calculate the date threshold using SQLite's datetime function
	dateThreshold := time.Now().AddDate(0, 0, -days).UTC().Format("2006-01-02")

	// Use strftime for date formatting in SQLite
	result := initializers.DB.Raw(`
		SELECT strftime('%Y-%m-%d', created_at) as date, COUNT(*) as count
		FROM transcriptions
		WHERE doctor_id = ?
		AND date(created_at) >= date(?)
		GROUP BY date
		ORDER BY count DESC
		LIMIT 5
	`, doctorUUID, dateThreshold).Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	return stats, nil
}

// SQLite-specific version of GetDashboardPatients for testing
func getDashboardPatientsForTest(doctorID string, days int) ([]models.Patient, error) {
	var patients []models.Patient

	// Return empty slice when days is 0
	if days == 0 {
		return patients, nil
	}

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID: %v", err)
	}

	// Calculate the date threshold using SQLite's datetime function
	dateThreshold := time.Now().AddDate(0, 0, -days).UTC().Format("2006-01-02 15:04:05")

	result := initializers.DB.Session(&gorm.Session{PrepareStmt: false}).
		Where("doctor_id = ? AND created_at >= ?", doctorUUID, dateThreshold).
		Order("created_at DESC").
		Find(&patients)

	if result.Error != nil {
		return nil, result.Error
	}

	return patients, nil
}
