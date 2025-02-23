package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Save original JWT_SECRET value
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	tests := []struct {
		name        string
		email       string
		secretKey   string
		shouldError bool
	}{
		{
			name:        "Valid token generation with environment secret",
			email:       "test@example.com",
			secretKey:   "test_secret_key",
			shouldError: false,
		},
		{
			name:        "Valid token generation with default secret",
			email:       "test@example.com",
			secretKey:   "",
			shouldError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			secretKey:   "test_secret_key",
			shouldError: false, // JWT spec doesn't require non-empty claims
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for the test
			os.Setenv("JWT_SECRET", tt.secretKey)

			// Generate token
			token, err := GenerateToken(tt.email)

			// Check error expectation
			if tt.shouldError {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Parse and verify the token
			claims := jwt.MapClaims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				secretKey := os.Getenv("JWT_SECRET")
				if secretKey == "" {
					secretKey = "default_secret_key_please_set_in_env"
				}
				return []byte(secretKey), nil
			})

			assert.NoError(t, err)
			assert.True(t, parsedToken.Valid)

			// Verify claims
			assert.Equal(t, tt.email, claims["email"])
			
			// Verify expiration time
			exp := time.Unix(int64(claims["exp"].(float64)), 0)
			expectedExp := time.Now().Add(time.Hour * 24)
			assert.WithinDuration(t, expectedExp, exp, time.Minute) // Allow 1 minute difference
		})
	}
}
