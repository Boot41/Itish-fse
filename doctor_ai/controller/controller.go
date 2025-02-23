package controller

import (
	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/middleware"
	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Works
func SignUp(c *gin.Context) {
	var user models.Doctor                          // creating a variable of type Doctor
	if err := c.ShouldBindJSON(&user); err != nil { // Binding the data coming from the request
		log.Println("Unable to Bind the data from frontend...")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input...",
		})
		return
	}
	log.Println("Bind Successfully from the frontend...")

	err := service.CreateDoctor(&user) // accessing the service layer and passing on the processed data to the service
	if err != nil {
		log.Println("Unable to create the user on the server side...")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to create user on the server side...",
		})
		return
	}
	log.Println("User created successfully on the server side...")
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

// Works
func Login(c *gin.Context) {
	var user models.Doctor                          // creating a variable of type User
	if err := c.ShouldBindJSON(&user); err != nil { // Binding the data coming from the request
		log.Println("Unable to Bind the data from frontend...")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input...",
		})
		return
	}
	log.Println("Bind Successfully from the frontend...")

	err := service.DoctorLogin(&user)
	if err != nil {
		log.Println("Unable to login the user on the server side...")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to login user on the server side...",
		})
		return
	}

	// Generate a token for the user
	token, err := middleware.GenerateToken(user.Email)
	if err != nil {
		log.Println("Error generating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate authentication token",
		})
		return
	}

	log.Println("User successfully logged in on the server side...")
	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
		"email":   user.Email,
		"token":   token,
	})
}

// Works
func GetProfile(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Extract email from request (assuming it’s sent in headers or body)
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Fetch doctor profile from the database
	profile, err := service.GetDoctorProfile(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// Works
func UpdateProfile(c *gin.Context) {
	var request struct {
		Email          string `json:"email"`
		Name           string `json:"name"`
		Specialization string `json:"specialization"`
		Phone          string `json:"phone"`
	}

	// Bind the JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body...")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Ensure email is provided
	if request.Email == "" {
		log.Println("Email is required to update the profile")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Call the service layer to update the doctor's profile
	err := service.UpdateDoctorProfile(request.Email, request.Name, request.Specialization, request.Phone)
	if err != nil {
		log.Println("Error updating the doctor profile:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Doctor profile updated successfully"})
}

// Works
// GetTranscriptions retrieves transcriptions for a doctor using the doctor email passed from the frontend.
func GetTranscriptions(c *gin.Context) {
	// Define a request struct to bind JSON payload
	var request struct {
		Email string `json:"email"`
	}

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Check if email is provided
	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Look up the doctor in the database.
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Handle pagination parameters.
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	// Call the service layer to get the transcriptions and patient names.
	formattedTranscriptions, patientNames, total, err := service.GetTranscriptions(doctor.ID, pageNum, limitNum)
	if err != nil {
		log.Println("Error retrieving transcriptions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve transcriptions"})
		return
	}

	// Return the results.
	c.JSON(http.StatusOK, gin.H{
		"transcripts": formattedTranscriptions,
		"patients":    patientNames,
		"total":       total,
		"page":        pageNum,
		"limit":       limitNum,
	})
}

// Works
func GetTranscriptionByID(c *gin.Context) {
	var request struct {
		TranscriptionID string `json:"transcription_id"`
	}

	// Bind the JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	transcriptionID, err := uuid.Parse(request.TranscriptionID)
	if err != nil {
		log.Println("Invalid transcription ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transcription ID"})
		return
	}

	transcription, err := service.GetTranscriptionByID(transcriptionID)
	if err != nil {
		log.Println("Error retrieving transcription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve transcription"})
		return
	}
	// log.Println("Retrieved transcription:", transcription)

	c.JSON(http.StatusOK, gin.H{
		"text":   transcription.Text,
		"report": transcription.Report,
	})
}

// Works
func UpdateTranscription(c *gin.Context) {
	// Get ID from URL parameter
	transcriptionIDStr := c.Param("id")
	transcriptionID, err := uuid.Parse(transcriptionIDStr)
	if err != nil {
		log.Println("Invalid transcription ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transcription ID"})
		return
	}

	var request struct {
		TranscriptionID string                 `json:"transcription_id"`
		UpdateData      map[string]interface{} `json:"update_data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	err = service.UpdateTranscription(transcriptionID, request.UpdateData)
	if err != nil {
		log.Println("Error updating transcription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transcription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transcription updated successfully"})
}

func DeleteTranscription(c *gin.Context) {

	transcriptionIDStr := c.Param("id")
	transcriptionID, err := uuid.Parse(transcriptionIDStr)
	if err != nil {
		log.Println("Invalid transcription ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transcription ID"})
		return
	}

	var request struct {
		TranscriptionID string `json:"transcription_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	err = service.DeleteTranscription(transcriptionID)
	if err != nil {
		log.Println("Error deleting transcription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete transcription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transcription deleted successfully"})
}

// Works
// DownloadTranscriptionRequest represents the JSON payload expected from the frontend.
type DownloadTranscriptionRequest struct {
	DoctorEmail string         `json:"email"`
	Patient     models.Patient `json:"patient"` // Must contain at least the patient ID
}

// Works
func DownloadTranscription(c *gin.Context) {
	// Bind JSON payload containing doctorEmail and patient object.
	var req DownloadTranscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Lookup the doctor by email.
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", req.DoctorEmail).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Retrieve the transcription for this doctor and patient.
	transcription, err := service.GetTranscriptionByPatient(doctor.ID, req.Patient.ID)
	if err != nil {
		log.Println("Error retrieving transcription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve transcription"})
		return
	}

	// Generate a PDF from the transcription report and patient details.
	pdfPath, err := service.GeneratePDF(transcription.Report, req.Patient)
	if err != nil {
		log.Println("Error generating PDF:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate PDF"})
		return
	}

	// Serve the generated PDF file.
	c.File(pdfPath)
}

// Works
// CreatePatientRequest represents the expected JSON payload from the frontend.
type CreatePatientRequest struct {
	DoctorEmail string         `json:"email"`
	Patient     models.Patient `json:"patient"`
	Audio       string         `json:"audio"` // audio URL as a string
}

// Works
// CreatePatientAndTranscription is the controller that handles the POST /patients endpoint.
func CreatePatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	transcription, err := service.CreatePatientAndTranscription(req.DoctorEmail, req.Patient, req.Audio)
	if err != nil {
		log.Println("Error creating patient and transcription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create patient/transcription"})
		return
	}

	c.JSON(http.StatusOK, transcription)
}

// Works
func GetPatients(c *gin.Context) {
	// Extract email from request
	var request struct {
		Email string `json:"email"`
	}

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Check if email is provided
	if request.Email == "" {
		log.Println("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	log.Printf("Fetching patients for email: %s", request.Email)

	// Fetch doctor by email
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Printf("Doctor not found for email %s: %v", request.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor not found"})
		return
	}

	log.Printf("Doctor found: %v", doctor)

	// Handle pagination
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	// Call the service to fetch patients
	patients, total, err := service.GetPatients(doctor.ID, pageNum, limitNum)
	if err != nil {
		log.Printf("Error retrieving patients for doctor %v: %v", doctor.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve patients"})
		return
	}

	log.Printf("Retrieved %d patients out of total %d", len(patients), total)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"total":    total,
		"page":     pageNum,
		"limit":    limitNum,
	})
}

// Works
func GetPatientByID(c *gin.Context) {
	// Extract email and patient ID from JSON request
	var request struct {
		Email     string `json:"email"`
		PatientID string `json:"patient_id"`
	}

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate email
	if request.Email == "" {
		log.Println("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Fetch doctor by email
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor not found"})
		return
	}

	// Validate patient ID
	patientID, err := uuid.Parse(request.PatientID)
	if err != nil {
		log.Println("Invalid patient ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid patient ID"})
		return
	}

	// Retrieve patient and transcription
	patient, transcription, err := service.GetPatientWithTranscription(doctor.ID, patientID)
	if err != nil {
		log.Println("Error retrieving patient data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve patient data"})
		return
	}

	// Return both patient and transcription data
	// Send only the fields needed by the frontend
	c.JSON(http.StatusOK, gin.H{
		"patient": gin.H{
			"name":   patient.Name,
			"age":    patient.Age,
			"gender": patient.Gender,
		},
		"transcription": transcription,
	})
}

// send patient_id along with email
// Works
func DeletePatient(c *gin.Context) {
	// Extract email and patient ID from JSON request
	var request struct {
		Email     string `json:"email"`
		PatientID string `json:"patient_id"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate email
	if request.Email == "" {
		log.Println("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Validate patient ID
	if request.PatientID == "" {
		log.Println("Patient ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Patient ID is required"})
		return
	}

	// Validate patient ID format
	parsedID, err := uuid.Parse(request.PatientID)
	if err != nil {
		log.Println("Invalid patient ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid patient ID"})
		return
	}

	// Call the service to delete the patient
	err = service.DeletePatient(parsedID)
	if err != nil {
		log.Println("Error deleting patient:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}

type DashboardRequest struct {
	Email string `json:"email"`
	Days  int    `json:"days"`
}

func GetDashboardTranscripts(c *gin.Context) {
	var request DashboardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Lookup the doctor in the database using the provided email.
	var doctor models.Doctor
	if err := initializers.DB.Model(&models.Doctor{}).Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Use days from request or default to 30
	daysNum := request.Days
	if daysNum < 1 {
		daysNum = 30
	}

	// Call the service layer to retrieve the dashboard transcripts.
	transcriptions, err := service.GetDashboardTranscripts(doctor.ID.String(), daysNum)
	if err != nil {
		log.Println("Error retrieving dashboard transcripts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve dashboard transcripts"})
		return
	}

	c.JSON(http.StatusOK, transcriptions)
}

func GetDashboardPatients(c *gin.Context) {
	var request DashboardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Lookup the doctor in the database using the provided email.
	var doctor models.Doctor
	if err := initializers.DB.Model(&models.Doctor{}).Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Use days from request or default to 30
	daysNum := request.Days
	if daysNum < 1 {
		daysNum = 30
	}

	// Retrieve the patients for the doctor using the service layer.
	patients, err := service.GetDashboardPatients(doctor.ID.String(), daysNum)
	if err != nil {
		log.Println("Error retrieving dashboard patients:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve dashboard patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

func GetDailyStatistics(c *gin.Context) {
	var request DashboardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Lookup the doctor by email.
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Use days from request or default to 30
	daysNum := request.Days
	if daysNum < 1 {
		daysNum = 30
	}

	// Call the service layer to retrieve daily statistics.
	stats, err := service.GetDailyStatistics(doctor.ID.String(), daysNum)
	if err != nil {
		log.Println("Error retrieving daily statistics:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve daily statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

type MonthlyStatsRequest struct {
	Email  string `json:"email"`
	Months int    `json:"months"`
}

func GetMonthlyStatistics(c *gin.Context) {
	var request MonthlyStatsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Lookup the doctor in the database using the provided email.
	var doctor models.Doctor
	if err := initializers.DB.Model(&models.Doctor{}).Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Use months from request or default to 12
	monthsNum := request.Months
	if monthsNum < 1 {
		monthsNum = 12
	}

	// Retrieve the monthly statistics for the doctor using the service layer.
	stats, err := service.GetMonthlyStatistics(doctor.ID.String(), monthsNum)
	if err != nil {
		log.Println("Error retrieving monthly statistics:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve monthly statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetBusiestDays(c *gin.Context) {
	var request DashboardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if request.Email == "" {
		log.Println("Doctor email not provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Doctor email is required"})
		return
	}

	// Lookup the doctor record using the provided email.
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Get the "days" query parameter (defaulting to 30).
	daysStr := c.DefaultQuery("days", "30")
	daysNum, err := strconv.Atoi(daysStr)
	if err != nil || daysNum < 1 {
		daysNum = 30
	}

	// Retrieve busiest days from the service layer.
	stats, err := service.GetBusiestDays(doctor.ID.String(), daysNum)
	if err != nil {
		log.Println("Error retrieving busiest days:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve busiest days"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Works
func UpdatePatient(c *gin.Context) {
	// Extract email and patient ID from JSON request
	var request struct {
		Email      string                 `json:"email"`
		PatientID  string                 `json:"patient_id"`
		UpdateData map[string]interface{} `json:"update_data"`
	}

	// Bind JSON request to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Validate email
	if request.Email == "" {
		log.Println("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	// Fetch doctor by email
	var doctor models.Doctor
	if err := initializers.DB.Where("email = ?", request.Email).First(&doctor).Error; err != nil {
		log.Println("Doctor not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Doctor not found"})
		return
	}

	// Validate patient ID
	patientID, err := uuid.Parse(request.PatientID)
	if err != nil {
		log.Println("Invalid patient ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid patient ID"})
		return
	}

	// Call the service to update the patient
	err = service.UpdatePatientByID(doctor.ID, patientID, request.UpdateData)
	if err != nil {
		log.Println("Error updating patient:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient updated successfully"})
}
