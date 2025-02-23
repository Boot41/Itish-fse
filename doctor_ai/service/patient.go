package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"log"
	"os"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// assemblyaiTranscribe sends the audio URL to AssemblyAI and returns the raw transcription text.
// For simplicity, this example assumes the transcription is ready.
func assemblyaiTranscribe(audioURL string) (string, error) {
	client := aai.NewClient(os.Getenv("ASSEMBLYAIKEY"))
	ctx := context.Background()

	transcript, err := client.Transcripts.TranscribeFromURL(ctx, audioURL, &aai.TranscriptOptionalParams{})
	if err != nil {
		return "", fmt.Errorf("assemblyai transcription error: %v", err)
	}

	if transcript.Text == nil {
		return "", fmt.Errorf("transcription text is empty")
	}
	return *transcript.Text, nil
}

// groqEnhance sends the raw transcription text to Groq to get a structured medical report.
func groqEnhance(rawText string) (string, error) {
	client := resty.New()

	requestBody := map[string]interface{}{
		"model": "llama3-8b-8192",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an AI specialized in generating well-structured medical reports. Ensure the report follows a professional format.",
			},
			{
				"role": "user",
				"content": fmt.Sprintf(
					"Format the following transcription into a structured medical report with these sections:\n"+
						"1. Patient Information\n"+
						"2. Patient History\n"+
						"3. Symptoms\n"+
						"4. Diagnosis\n"+
						"5. Treatment Plan\n"+
						"6. Recommendations\n\nTranscription:\n%s", rawText),
			},
		},
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+os.Getenv("GROQAPIKEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post("https://api.groq.com/openai/v1/chat/completions")
	if err != nil {
		return "", fmt.Errorf("groq request error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal groq response: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected response from Groq API")
	}

	content, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract enhanced content from Groq response")
	}

	return content, nil
}

func CreatePatientAndTranscription(doctorEmail string, patientData models.Patient, audioURL string) (*models.Transcription, error) {
	if err := validatePatient(&patientData); err != nil {
		log.Println("Validation error:", err)
		return nil, err
	}
	log.Println("Patient validated successfully:", patientData.ID)

	// Step 1: Retrieve doctor by email
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", doctorEmail).First(&doctor).Error; err != nil {
		log.Println("Error retrieving doctor:", err)
		return nil, fmt.Errorf("doctor not found: %v", err)
	}
	log.Println("Doctor found:", doctor.ID)

	// Step 2: Create the patient record and link it to the doctor
	patientData.DoctorID = doctor.ID
	if err := initializers.DB.Create(&patientData).Error; err != nil {
		log.Println("Error creating patient:", err)
		return nil, fmt.Errorf("failed to create patient: %v", err)
	}
	log.Println("Patient created successfully:", patientData.ID)

	// Step 3: Transcribe the audio using AssemblyAI
	rawTranscript, err := assemblyaiTranscribe(audioURL)
	if err != nil {
		log.Println("Error transcribing audio:", err)
		return nil, fmt.Errorf("failed to transcribe audio: %v", err)
	}
	log.Println("Transcription created successfully:", patientData.ID)

	// Step 4: Enhance the transcription using Groq
	enhancedTranscript, err := groqEnhance(rawTranscript)
	if err != nil {
		log.Println("Error enhancing transcription:", err)
		return nil, fmt.Errorf("failed to enhance transcription: %v", err)
	}
	log.Println("Enhanced transcription created successfully:", patientData.ID)

	// Step 5: Create the transcription record in the database
	newTranscription := models.Transcription{
		ID:        uuid.New(),
		DoctorID:  doctor.ID,
		PatientID: patientData.ID,
		Text:      rawTranscript,
		Report:    enhancedTranscript,
		CreatedAt: time.Now(),
	}

	if err := initializers.DB.Create(&newTranscription).Error; err != nil {
		log.Println("Error creating transcription record:", err)
		return nil, fmt.Errorf("failed to create transcription record: %v", err)
	}
	log.Println("Transcription record created successfully:", newTranscription.ID)
	return &newTranscription, nil
}

// Example validation function (You can modify this)
func validatePatient(patient *models.Patient) error {
	if patient.Name == "" {
		return errors.New("patient name is required")
	}
	if patient.Age <= 0 {
		return errors.New("patient age must be greater than zero")
	}
	return nil
}

func GetPatients(doctorID uuid.UUID, page, limit int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	// Count total patients for the doctor
	if err := initializers.DB.Model(&models.Patient{}).Where("doctor_id = ?", doctorID).Count(&total).Error; err != nil {
		log.Println("Error counting patients:", err)
		return nil, 0, errors.New("failed to count patients")
	}

	// Retrieve patients with pagination
	offset := (page - 1) * limit
	if err := initializers.DB.Where("doctor_id = ?", doctorID).Limit(limit).Offset(offset).Find(&patients).Error; err != nil {
		log.Println("Error fetching patients:", err)
		return nil, 0, errors.New("failed to retrieve patients")
	}

	return patients, total, nil
}

func GetPatientWithTranscription(doctorID uuid.UUID, patientID uuid.UUID) (*models.Patient, *models.Transcription, error) {
	var patient models.Patient
	var transcription models.Transcription

	// Retrieve patient only if it belongs to the doctor
	if err := initializers.DB.Where("id = ? AND doctor_id = ?", patientID, doctorID).First(&patient).Error; err != nil {
		log.Println("Error fetching patient:", err)
		return nil, nil, errors.New("patient not found or does not belong to the doctor")
	}

	// Retrieve transcription associated with the patient
	if err := initializers.DB.Where("patient_id = ?", patientID).First(&transcription).Error; err != nil {
		log.Println("Error fetching transcription:", err)
		return &patient, nil, nil // Patient exists, but no transcription found
	}

	return &patient, &transcription, nil
}

func UpdatePatientByID(doctorID uuid.UUID, patientID uuid.UUID, updateData map[string]interface{}) error {
	var patient models.Patient

	// Check if patient exists and belongs to the doctor
	if err := initializers.DB.Where("id = ? AND doctor_id = ?", patientID, doctorID).First(&patient).Error; err != nil {
		log.Println("Patient not found or does not belong to the doctor:", err)
		return errors.New("patient not found or unauthorized")
	}

	// Update only allowed fields
	allowedFields := map[string]bool{
		"name":   true,
		"age":    true,
		"gender": true,
	}

	for key := range updateData {
		if !allowedFields[key] {
			delete(updateData, key) // Remove disallowed fields
		}
	}

	// Update the patient record
	if err := initializers.DB.Model(&patient).Updates(updateData).Error; err != nil {
		log.Println("Error updating patient:", err)
		return errors.New("failed to update patient")
	}

	log.Println("Patient updated successfully:", patientID)
	return nil
}

// Works
func DeletePatient(patientID uuid.UUID) error {
	var patient models.Patient

	// Check if patient exists
	if err := initializers.DB.Where("id = ?", patientID).First(&patient).Error; err != nil {
		log.Println("Patient not found:", err)
		return errors.New("patient not found")
	}

	// Delete transcriptions first (to maintain referential integrity)
	if err := initializers.DB.Where("patient_id = ?", patientID).Delete(&models.Transcription{}).Error; err != nil {
		log.Println("Error deleting transcriptions:", err)
		return errors.New("failed to delete transcriptions")
	}

	// Delete the patient
	if err := initializers.DB.Delete(&patient).Error; err != nil {
		log.Println("Failed to delete patient:", err)
		return errors.New("failed to delete patient")
	}

	return nil
}
