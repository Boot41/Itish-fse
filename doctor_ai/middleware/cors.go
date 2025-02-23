package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set allowed origins (your frontend services)
		allowedOrigins := []string{
			"https://dreamnote-react.onrender.com", // Production React frontend
			"http://localhost:5173",                // Local Vite development
			"http://localhost:3000",                // Docker frontend
		}

		// Check the Origin header from the request
		origin := c.Request.Header.Get("Origin")
		isAllowed := false

		// Validate if the origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin ||
				(allowedOrigin == "http://localhost:5173" &&
					(origin == "http://localhost:5173" ||
						origin == "http://localhost:5173/" ||
						origin == "http://localhost:5173/*")) ||
				(allowedOrigin == "http://localhost:3000" &&
					(origin == "http://localhost:3000" ||
						origin == "http://localhost:3000/" ||
						origin == "http://localhost:3000/*")) {
				isAllowed = true
				break
			}
		}

		// Set CORS headers based on origin validation
		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Fallback for development or unrecognized origins
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		// Enable credentials
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Common CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Requested-With, Accept")

		// Handle preflight requests (OPTIONS method)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}
