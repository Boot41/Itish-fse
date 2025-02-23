package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUpdateProfile_InvalidJSON verifies that an invalid JSON payload returns a 400 Bad Request.
func TestUpdateProfile_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	UpdateProfile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateProfile_EmailMissing verifies that if the email is missing, UpdateProfile returns a 400.
func TestUpdateProfile_EmailMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Prepare payload with missing email.
	payload := map[string]string{
		"email":          "",
		"name":           "Dr. Test",
		"specialization": "Cardiology",
		"phone":          "1234567890",
	}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	UpdateProfile(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Optionally, you can check the response message.
}

// TestUpdateProfile_ServiceFailure patches service.UpdateDoctorProfile to return an error.
func TestUpdateProfile_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := map[string]string{
		"email":          "test@example.com",
		"name":           "Dr. Test",
		"specialization": "Cardiology",
		"phone":          "1234567890",
	}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.UpdateDoctorProfile to simulate an error.
	patches := gomonkey.ApplyFunc(service.UpdateDoctorProfile, func(email, name, specialization, phone string) error {
		return errors.New("update error")
	})
	defer patches.Reset()

	UpdateProfile(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "update error", resp["message"])
}

// TestUpdateProfile_Success patches service.UpdateDoctorProfile to succeed.
func TestUpdateProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := map[string]string{
		"email":          "test@example.com",
		"name":           "Dr. Test",
		"specialization": "Cardiology",
		"phone":          "1234567890",
	}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/updateprofile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.UpdateDoctorProfile to simulate success.
	patches := gomonkey.ApplyFunc(service.UpdateDoctorProfile, func(email, name, specialization, phone string) error {
		return nil
	})
	defer patches.Reset()

	UpdateProfile(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Doctor profile updated successfully", resp["message"])
}
