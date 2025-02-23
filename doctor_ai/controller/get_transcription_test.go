package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/models"
	"itish41/doctor_ai_assistant/service"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTranscriptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocking DB queries properly
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "Where", func(_ *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
		return &gorm.DB{}
	})
	defer patches.Reset()

	patches2 := gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "First", func(_ *gorm.DB, dest interface{}) *gorm.DB {
		doctor, ok := dest.(*models.Doctor)
		if ok {
			doctor.ID = uuid.New() // Ensure it's UUID
			doctor.Email = "test@example.com"
		}
		return &gorm.DB{}
	})
	defer patches2.Reset()

	// Mock GetTranscriptions response
	patches3 := gomonkey.ApplyFunc(service.GetTranscriptions, func(doctorID uuid.UUID, page, limit int) ([]map[string]interface{}, []string, int64, error) {
		return []map[string]interface{}{
			{"text": "Test Transcript", "report": "Sample Report"},
		}, []string{"John Doe"}, 1, nil
	})
	defer patches3.Reset()

	router := gin.Default()
	router.POST("/get_transcriptions", GetTranscriptions)
	requestBody := `{"email": "test@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcriptions", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check total count
	totalFloat, ok := response["total"].(float64)
	assert.True(t, ok, "total should be a number")
	assert.Equal(t, int64(1), int64(totalFloat))
}

func TestGetTranscriptions_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/get_transcriptions", GetTranscriptions)
	requestBody := `{"email": ""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcriptions", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func setupMockDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect mock database")
	}
	initializers.DB = db
}

func TestGetTranscriptions_DoctorNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setupMockDB() // Set up mock DB connection

	router := gin.Default()
	router.POST("/get_transcriptions", GetTranscriptions)
	requestBody := `{"email": "nonexistent@example.com"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get_transcriptions", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
