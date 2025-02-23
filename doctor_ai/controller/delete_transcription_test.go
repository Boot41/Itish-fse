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

func TestDeleteTranscription_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock service.DeleteTranscription to return nil (success)
	patches := gomonkey.ApplyFuncReturn(service.DeleteTranscription, nil)
	defer patches.Reset()

	router := gin.Default()
	router.DELETE("/delete_transcription/:id", DeleteTranscription)

	transcriptionID := uuid.New().String()
	requestBody := `{
		"transcription_id": "` + transcriptionID + `"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete_transcription/"+transcriptionID, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transcription deleted successfully", response["message"])
}

func TestDeleteTranscription_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.DELETE("/delete_transcription/:id", DeleteTranscription)

	requestBody := `{
		"transcription_id": "invalid-uuid"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete_transcription/invalid-uuid", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTranscription_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock service.DeleteTranscription to return an error
	patches := gomonkey.ApplyFuncReturn(service.DeleteTranscription, assert.AnError)
	defer patches.Reset()

	router := gin.Default()
	router.DELETE("/delete_transcription/:id", DeleteTranscription)

	transcriptionID := uuid.New().String()
	requestBody := `{
		"transcription_id": "` + transcriptionID + `"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete_transcription/"+transcriptionID, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
