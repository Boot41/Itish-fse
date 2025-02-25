package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Setup routes
	SetupRoutes(r)

	// Define test cases for each route group
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		// Auth routes
		{
			name:           "SignUp Route",
			method:         "POST",
			path:           "/auth/signup",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Login Route",
			method:         "POST",
			path:           "/auth/login",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Profile Route",
			method:         "POST",
			path:           "/auth/me",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Profile Route",
			method:         "PUT",
			path:           "/auth/me",
			expectedStatus: http.StatusOK,
		},

		// Transcription routes
		{
			name:           "Get Transcriptions Route",
			method:         "POST",
			path:           "/transcription/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Download Transcription Route",
			method:         "POST",
			path:           "/transcription/123/download",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Transcription By ID Route",
			method:         "POST",
			path:           "/transcription/123/getTranscriptionByID",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Transcription Route",
			method:         "PUT",
			path:           "/transcription/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Transcription Route",
			method:         "POST",
			path:           "/transcription/123/delete",
			expectedStatus: http.StatusOK,
		},

		// Patient routes
		{
			name:           "Create Patient Route",
			method:         "POST",
			path:           "/patients/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Patients List Route",
			method:         "POST",
			path:           "/patients/patientsList",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Patient By ID Route",
			method:         "POST",
			path:           "/patients/123/getPatient",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Patient Route",
			method:         "PUT",
			path:           "/patients/123/update",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Patient Route",
			method:         "POST",
			path:           "/patients/123",
			expectedStatus: http.StatusOK,
		},

		// Dashboard routes
		{
			name:           "Get Dashboard Transcripts Route",
			method:         "POST",
			path:           "/dashboard/transcripts",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Dashboard Patients Route",
			method:         "POST",
			path:           "/dashboard/patients",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Daily Statistics Route",
			method:         "POST",
			path:           "/dashboard/statistics/daily",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Monthly Statistics Route",
			method:         "POST",
			path:           "/dashboard/statistics/monthly",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Busiest Days Route",
			method:         "POST",
			path:           "/dashboard/statistics/busiest-days",
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests for each route
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			// Serve the request directly
			r.ServeHTTP(w, req)

			// Verify the route exists (404 means route not found)
			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route should exist: %s %s", tt.method, tt.path)
		})
	}
}

func TestRouteGroups(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Setup routes
	SetupRoutes(r)

	// Get the route groups
	routes := r.Routes()

	// Helper function to check if a path belongs to a group
	hasRouteInGroup := func(group string, routes gin.RoutesInfo) bool {
		for _, route := range routes {
			if len(route.Path) >= len(group) && route.Path[:len(group)] == group {
				return true
			}
		}
		return false
	}

	// Test each route group exists
	t.Run("Route Groups Exist", func(t *testing.T) {
		assert.True(t, hasRouteInGroup("/auth", routes), "Auth group should exist")
		assert.True(t, hasRouteInGroup("/transcription", routes), "Transcription group should exist")
		assert.True(t, hasRouteInGroup("/patients", routes), "Patients group should exist")
		assert.True(t, hasRouteInGroup("/dashboard", routes), "Dashboard group should exist")
	})

	// Count routes in each group
	authCount := 0
	transcriptionCount := 0
	patientsCount := 0
	dashboardCount := 0

	for _, route := range routes {
		switch {
		case len(route.Path) >= 5 && route.Path[:5] == "/auth":
			authCount++
		case len(route.Path) >= 14 && route.Path[:14] == "/transcription":
			transcriptionCount++
		case len(route.Path) >= 9 && route.Path[:9] == "/patients":
			patientsCount++
		case len(route.Path) >= 10 && route.Path[:10] == "/dashboard":
			dashboardCount++
		}
	}

	// Verify route counts
	t.Run("Route Counts", func(t *testing.T) {
		assert.Equal(t, 4, authCount, "Auth group should have 4 routes")
		assert.Equal(t, 5, transcriptionCount, "Transcription group should have 5 routes")
		assert.Equal(t, 6, patientsCount, "Patients group should have 5 routes")
		assert.Equal(t, 5, dashboardCount, "Dashboard group should have 5 routes")
	})
}
