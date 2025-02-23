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

// TestSignUp_InvalidJSON verifies that SignUp returns a 400 Bad Request
// when an invalid JSON payload is sent.
func TestSignUp_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	SignUp(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSignUp_ServiceFailure simulates a failure in service.CreateDoctor,
// expecting the SignUp handler to return a 500 Internal Server Error.
func TestSignUp_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare valid JSON payload for a Doctor.
	user := models.Doctor{Email: "test@example.com"}
	jsonData, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.CreateDoctor to simulate an error.
	patches := gomonkey.ApplyFunc(service.CreateDoctor, func(u *models.Doctor) error {
		return errors.New("service error")
	})
	defer patches.Reset()

	SignUp(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestSignUp_Success simulates a successful service.CreateDoctor call,
// expecting the SignUp handler to return a 200 OK.
func TestSignUp_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare valid JSON payload for a Doctor.
	user := models.Doctor{Email: "test@example.com"}
	jsonData, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Patch service.CreateDoctor to simulate success.
	patches := gomonkey.ApplyFunc(service.CreateDoctor, func(u *models.Doctor) error {
		return nil
	})
	defer patches.Reset()

	SignUp(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
