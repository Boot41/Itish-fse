package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTranscriptionByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	transcriptionID := uuid.New()

	// Mock service function
	patches := gomonkey.ApplyFunc(service.GetTranscriptionByID, func(id uuid.UUID) (*models.Transcription, error) {
		return &models.Transcription{
			Text:   "Sample transcription text",
			Report: "Sample report",
		}, nil
	})
	defer patches.Reset()

	router := gin.Default()
	router.POST("/get_transcription", GetTranscriptionByID)
	requestBody := `{"transcription_id": "` + transcriptionID.String() + `"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcription", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Sample transcription text", response["text"])
	assert.Equal(t, "Sample report", response["report"])
}

func TestGetTranscriptionByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/get_transcription", GetTranscriptionByID)
	requestBody := `{"transcription_id": "invalid-uuid"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcription", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTranscriptionByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	transcriptionID := uuid.New()

	// Mock service function to return an error
	patches := gomonkey.ApplyFunc(service.GetTranscriptionByID, func(id uuid.UUID) (*models.Transcription, error) {
		return nil, assert.AnError
	})
	defer patches.Reset()

	router := gin.Default()
	router.POST("/get_transcription", GetTranscriptionByID)
	requestBody := `{"transcription_id": "` + transcriptionID.String() + `"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcription", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
