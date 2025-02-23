package initializers

import (
	"testing"
)

func TestConnectDB(t *testing.T) {
	// Set the Supabase connection string as the environment variable.
	// Make sure to replace the values with those provided by your Supabase dashboard.
	testDSN := "postgresql://postgres.mwaowvtmopqrbknfojxs:Itish@Kiit@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"
	t.Setenv("DIRECT_URL", testDSN)

	// Call the ConnectDB function from the original file.
	err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify that the global DB variable is not nil.
	if DB == nil {
		t.Fatal("Expected DB to be initialized, but got nil")
	}
}
