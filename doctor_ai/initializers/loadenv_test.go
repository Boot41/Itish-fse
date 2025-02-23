package initializers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv_Success(t *testing.T) {
	// Create a temporary directory for the test.
	tempDir := t.TempDir()

	// Write a .env file with some dummy content.
	envContent := "FOO=bar\n"
	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	// Save current working directory.
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	// Change to the temporary directory so that godotenv.Load() finds the .env file.
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	// Restore original working directory when the test is done.
	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	// Call LoadEnv and expect no error.
	if err := LoadEnv(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify that the environment variable from the .env file was loaded.
	if v := os.Getenv("FOO"); v != "bar" {
		t.Errorf("expected FOO=bar, got FOO=%s", v)
	}
}

func TestLoadEnv_Failure(t *testing.T) {
	// Create a temporary directory with no .env file.
	tempDir := t.TempDir()

	// Save current working directory.
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	// Change to the temporary directory.
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	// Restore the original working directory when the test finishes.
	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	// Confirm that no .env file exists.
	if _, err := os.Stat(filepath.Join(tempDir, ".env")); err == nil {
		t.Fatal("expected no .env file in the temporary directory")
	}

	// Call LoadEnv, which should fail because no .env file exists.
	err = LoadEnv()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != "env not loading" {
		t.Errorf("expected error message 'env not loading', got: %v", err)
	}
}
