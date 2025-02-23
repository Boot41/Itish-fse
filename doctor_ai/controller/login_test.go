package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"itish41/doctor_ai_assistant/middleware"
	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestLogin_InvalidJSON tests that an invalid JSON payload returns 400.
func TestLogin_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestLogin_ServiceFailure tests that if service.DoctorLogin returns an error, Login returns 500.
func TestLogin_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare a valid JSON payload for a Doctor.
	user := models.Doctor{Email: "test@example.com"}
	jsonData, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.DoctorLogin to simulate an error.
	patches := gomonkey.ApplyFunc(service.DoctorLogin, func(u *models.Doctor) error {
		return errors.New("login error")
	})
	defer patches.Reset()

	Login(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestLogin_TokenFailure tests that if middleware.GenerateToken fails (after successful login),
// Login returns a 500 error.
func TestLogin_TokenFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := models.Doctor{Email: "test@example.com"}
	jsonData, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.DoctorLogin to succeed.
	patches1 := gomonkey.ApplyFunc(service.DoctorLogin, func(u *models.Doctor) error {
		return nil
	})
	defer patches1.Reset()

	// Patch middleware.GenerateToken to simulate a token generation error.
	patches2 := gomonkey.ApplyFunc(middleware.GenerateToken, func(email string) (string, error) {
		return "", errors.New("token error")
	})
	defer patches2.Reset()

	Login(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestLogin_Success tests the successful login scenario.
func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := models.Doctor{Email: "test@example.com"}
	jsonData, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.DoctorLogin to succeed.
	patches1 := gomonkey.ApplyFunc(service.DoctorLogin, func(u *models.Doctor) error {
		return nil
	})
	defer patches1.Reset()

	// Patch middleware.GenerateToken to return a dummy token.
	patches2 := gomonkey.ApplyFunc(middleware.GenerateToken, func(email string) (string, error) {
		return "dummy-token", nil
	})
	defer patches2.Reset()

	Login(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Optionally, verify the response contains the expected token and email.
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "User login successfully", resp["message"])
	assert.Equal(t, "test@example.com", resp["email"])
	assert.Equal(t, "dummy-token", resp["token"])
}
