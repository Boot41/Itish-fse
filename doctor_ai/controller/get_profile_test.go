package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetProfile_InvalidJSON verifies that an invalid JSON payload returns 400.
func TestGetProfile_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GetProfile(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetProfile_EmailMissing verifies that an empty email returns 400.
func TestGetProfile_EmailMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := map[string]string{"email": ""}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	GetProfile(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Email is required", resp["message"])
}

// TestGetProfile_ServiceFailure patches service.GetDoctorProfile to simulate an error.
func TestGetProfile_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := map[string]string{"email": "test@example.com"}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.GetDoctorProfile with a function that returns a nil pointer and an error.
	patches := gomonkey.ApplyFunc(service.GetDoctorProfile, func(email string) (*models.Doctor, error) {
		return nil, errors.New("not found")
	})
	defer patches.Reset()

	GetProfile(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Failed to retrieve profile", resp["message"])
}

// TestGetProfile_Success patches service.GetDoctorProfile to return a valid profile.
func TestGetProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := map[string]string{"email": "test@example.com"}
	jsonData, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Create a dummy profile.
	dummyProfile := &models.Doctor{
		Email: "test@example.com",
		// Populate additional fields as needed.
	}

	// Patch service.GetDoctorProfile to return the dummy profile.
	patches := gomonkey.ApplyFunc(service.GetDoctorProfile, func(email string) (*models.Doctor, error) {
		return dummyProfile, nil
	})
	defer patches.Reset()

	GetProfile(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Doctor
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, dummyProfile.Email, resp.Email)
}
