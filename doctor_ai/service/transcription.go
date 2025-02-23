package service

import (
	"errors"
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// func CreateTranscription(doctorID uuid.UUID, file *multipart.FileHeader, patientID *uuid.UUID) (*models.Transcription, error) {
// 	// Validate input
// 	if doctorID == uuid.Nil {
// 		return nil, errors.New("invalid doctor ID")
// 	}

// 	// Generate unique filename
// 	filename := fmt.Sprintf("%s_%s%s", doctorID.String(), time.Now().Format("20060102_150405"), filepath.Ext(file.Filename))
// 	filepath := filepath.Join("uploads", "transcriptions", filename)

// 	// Ensure uploads directory exists
// 	os.MkdirAll(filepath.Dir(filepath), os.ModePerm)

// 	// Save file
// 	if err := SaveUploadedFile(file, filepath); err != nil {
// 		log.Println("Error saving transcription file:", err)
// 		return nil, errors.New("failed to save transcription file")
// 	}

// 	// Create transcription record
// 	transcription := &models.Transcription{
// 		DoctorID:         doctorID,
// 		PatientID:        patientID,
// 		Filepath:         filepath,
// 		OriginalFilename: file.Filename,
// 		CreatedAt:        time.Now(),
// 	}

// 	result := initializers.DB.Create(transcription)
// 	if result.Error != nil {
// 		log.Println("Error creating transcription record:", result.Error)
// 		os.Remove(filepath) // Remove file if DB record fails
// 		return nil, errors.New("failed to create transcription record")
// 	}

// 	log.Println("Transcription created successfully:", transcription.ID)
// 	return transcription, nil
// }

// // func SaveUploadedFile(file *multipart.FileHeader, filepath string) error {
// // 	src, err := file.Open()
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer src.Close()

// // 	dst, err := os.Create(filepath)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer dst.Close()

// // 	_, err = dst.Write(src.Bytes())
// // 	return err
// // }

// GetTranscriptions retrieves the list of transcriptions for a given doctor with pagination.
func GetTranscriptions(doctorID uuid.UUID, page, limit int) ([]map[string]interface{}, []string, int64, error) {
	var transcriptions []models.Transcription
	var total int64

	// Calculate the offset based on page and limit.
	offset := (page - 1) * limit

	// Count the total number of transcriptions for the doctor.
	if err := initializers.DB.
		Model(&models.Transcription{}).
		Where("doctor_id = ?", doctorID).
		Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	// Retrieve the transcriptions with pagination.
	if err := initializers.DB.
		Where("doctor_id = ?", doctorID).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&transcriptions).Error; err != nil {
		return nil, nil, 0, err
	}

	// Prepare patient names list and formatted transcriptions
	patientNames := make([]string, len(transcriptions))
	formattedTranscriptions := make([]map[string]interface{}, len(transcriptions))

	for i, t := range transcriptions {
		// Get patient name
		var patient models.Patient
		if err := initializers.DB.
			Where("id = ?", t.PatientID).
			First(&patient).Error; err != nil {
			patientNames[i] = "Unknown Patient"
		} else {
			patientNames[i] = patient.Name
		}

		// Format transcription
		formattedTranscriptions[i] = map[string]interface{}{
			"id":            t.ID,
			"patientId":     t.PatientID,
			"timestamp":     t.CreatedAt,
			"text":          t.Text,
			"report":        t.Report,
			"transcriptUrl": fmt.Sprintf("/transcription/%s/download", t.ID),
		}
	}

	return formattedTranscriptions, patientNames, total, nil
}

// Works
func GetTranscriptionByID(transcriptionID uuid.UUID) (*models.Transcription, error) {
	var transcription models.Transcription
	result := initializers.DB.Where("id = ?", transcriptionID).First(&transcription)
	if result.Error != nil {
		log.Println("Error retrieving transcription:", result.Error)
		return nil, errors.New("transcription not found")
	}
	// log.Println("Retrieved transcription:", transcription)

	return &transcription, nil
}

func UpdateTranscription(transcriptionID uuid.UUID, updateData map[string]interface{}) error {
	result := initializers.DB.Model(&models.Transcription{}).
		Where("id = ?", transcriptionID).
		Updates(updateData)

	if result.Error != nil {
		log.Println("Error updating transcription:", result.Error)
		return errors.New("failed to update transcription")
	}

	log.Println("Transcription updated successfully:", transcriptionID)
	return nil
}

func DeleteTranscription(transcriptionID uuid.UUID) error {
	var transcription models.Transcription

	// First, find the transcription to get its filepath
	result := initializers.DB.Where("id = ?", transcriptionID).First(&transcription)
	if result.Error != nil {
		log.Println("Error finding transcription:", result.Error)
		return errors.New("transcription not found")
	}

	// // Delete the file
	// if err := os.Remove(transcription.Filepath); err != nil && !os.IsNotExist(err) {
	// 	log.Println("Error deleting transcription file:", err)
	// 	return errors.New("failed to delete transcription file")
	// }

	// Delete the database record
	result = initializers.DB.Delete(&transcription)
	if result.Error != nil {
		log.Println("Error deleting transcription record:", result.Error)
		return errors.New("failed to delete transcription record")
	}

	log.Println("Transcription deleted successfully:", transcriptionID)
	return nil
}

// GetTranscriptionByPatient retrieves the latest transcription for a given doctor and patient.
func GetTranscriptionByPatient(doctorID, patientID uuid.UUID) (*models.Transcription, error) {
	var transcription models.Transcription
	err := initializers.DB.
		Where("doctor_id = ? AND patient_id = ?", doctorID, patientID).
		Order("created_at desc").
		First(&transcription).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transcription: %v", err)
	}
	return &transcription, nil
}

// GeneratePDF creates a PDF document from the transcription report and patient details.
// It returns the file path to the generated PDF.
func GeneratePDF(report string, patient models.Patient) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add Header
	pdf.Cell(40, 10, "Medical Transcription Report")
	pdf.Ln(12)

	// Add Patient Information
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Patient ID: %s", patient.ID))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Patient Name: %s", patient.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Age: %d", patient.Age))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Gender: %s", patient.Gender))
	pdf.Ln(12)

	// Add Report Title
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Transcription Report:")
	pdf.Ln(10)

	// Add Report Content using MultiCell for long text
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 7, report, "", "L", false)

	// Generate file name with timestamp
	fileName := fmt.Sprintf("transcription_%d.pdf", time.Now().Unix())
	filePath := "./pdfs/" + fileName

	// Ensure the directory exists
	if err := os.MkdirAll("./pdfs", os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create pdf directory: %v", err)
	}

	// Save the PDF to file
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", fmt.Errorf("failed to save pdf: %v", err)
	}

	return filePath, nil
}

type Patient struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Age       int
	Gender    string
	Email     string
	DoctorID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	// Add other fields as needed
}

func GetPatientByID(id string) (*Patient, error) {
	var patient Patient
	result := initializers.DB.Where("id = ?", id).First(&patient)
	return &patient, result.Error
}
