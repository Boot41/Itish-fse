package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTranscription_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock service.UpdateTranscription
	patches := gomonkey.ApplyFuncReturn(service.UpdateTranscription, nil)
	defer patches.Reset()

	router := gin.Default()
	router.PUT("/update_transcription/:id", UpdateTranscription)

	transcriptionID := uuid.New().String()
	requestBody := `{
		"transcription_id": "` + transcriptionID + `",
		"update_data": {
			"text": "Updated text",
			"report": "Updated report"
		}
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update_transcription/"+transcriptionID, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transcription updated successfully", response["message"])
}

func TestUpdateTranscription_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.PUT("/update_transcription/:id", UpdateTranscription)

	requestBody := `{
		"transcription_id": "invalid-uuid",
		"update_data": {
			"text": "Updated text"
		}
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update_transcription/invalid-uuid", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTranscription_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock service.UpdateTranscription to return an error
	patches := gomonkey.ApplyFuncReturn(service.UpdateTranscription, assert.AnError)
	defer patches.Reset()

	router := gin.Default()
	router.PUT("/update_transcription/:id", UpdateTranscription)

	transcriptionID := uuid.New().String()
	requestBody := `{
		"transcription_id": "` + transcriptionID + `",
		"update_data": {
			"text": "Updated text"
		}
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update_transcription/"+transcriptionID, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
