package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		origin         string
		expectedOrigin string
		method         string
	}{
		{
			name:           "Production frontend origin",
			origin:         "https://dreamnote-react.onrender.com",
			expectedOrigin: "https://dreamnote-react.onrender.com",
			method:         "GET",
		},
		{
			name:           "Local Vite development origin",
			origin:         "http://localhost:5173",
			expectedOrigin: "http://localhost:5173",
			method:         "POST",
		},
		{
			name:           "Docker frontend origin",
			origin:         "http://localhost:3000",
			expectedOrigin: "http://localhost:3000",
			method:         "OPTIONS",
		},
		{
			name:           "Unknown origin",
			origin:         "http://unknown-origin.com",
			expectedOrigin: "http://localhost:3000", // Fallback origin
			method:         "GET",
		},
		{
			name:           "Local Vite with trailing slash",
			origin:         "http://localhost:5173/",
			expectedOrigin: "http://localhost:5173/",
			method:         "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			router.Use(CORSMiddleware())

			// Add a test endpoint
			router.Any("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create a test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			// Serve the request
			router.ServeHTTP(w, req)

			// Check CORS headers
			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization, Origin, X-Requested-With, Accept", w.Header().Get("Access-Control-Allow-Headers"))

			// For OPTIONS requests (preflight), check status code
			if tt.method == "OPTIONS" {
				assert.Equal(t, http.StatusNoContent, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}
