package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"itish41/doctor_ai_assistant/service"
)

func TestDeletePatient_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the service.DeletePatient function to return nil (success)
	patches := gomonkey.ApplyFunc(service.DeletePatient, func(patientID uuid.UUID) error {
		return nil
	})
	defer patches.Reset()

	router := gin.Default()
	router.POST("/delete_patient", DeletePatient)

	patientID := uuid.New().String()

	requestBody, _ := json.Marshal(map[string]string{
		"email":      "testdoctor@example.com",
		"patient_id": patientID,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete_patient", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Patient deleted successfully")
}

func TestDeletePatient_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/delete_patient", DeletePatient)

	requestBody, _ := json.Marshal(map[string]string{
		"email":      "testdoctor@example.com",
		"patient_id": "invalid-uuid",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete_patient", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid patient ID")
}

func TestDeletePatient_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock the service.DeletePatient function to return an error
	patches := gomonkey.ApplyFunc(service.DeletePatient, func(patientID uuid.UUID) error {
		return errors.New("database error")
	})
	defer patches.Reset()

	router := gin.Default()
	router.POST("/delete_patient", DeletePatient)

	patientID := uuid.New().String()

	requestBody, _ := json.Marshal(map[string]string{
		"email":      "testdoctor@example.com",
		"patient_id": patientID,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete_patient", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete patient")
}
