package initializers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	// 1. Create a temporary directory that will serve as our project root.
	tempDir := t.TempDir()

	// 2. Create the directory structure "db/migrations" inside tempDir.
	migrationsDir := filepath.Join(tempDir, "db", "migrations")
	err := os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err, "failed to create migrations directory")

	// 3. Create dummy migration files.
	//    Create an "up" migration that creates a table "dummy"
	//    and a corresponding "down" migration that drops it.
	upFilename := filepath.Join(migrationsDir, "0001_dummy.up.sql")
	downFilename := filepath.Join(migrationsDir, "0001_dummy.down.sql")

	upContent := []byte("CREATE TABLE IF NOT EXISTS dummy (id SERIAL PRIMARY KEY);")
	downContent := []byte("DROP TABLE IF EXISTS dummy;")

	err = os.WriteFile(upFilename, upContent, 0644)
	require.NoError(t, err, "failed to write up migration file")
	err = os.WriteFile(downFilename, downContent, 0644)
	require.NoError(t, err, "failed to write down migration file")

	// 4. Change the working directory to tempDir so that "file://db/migrations" resolves correctly.
	origDir, err := os.Getwd()
	require.NoError(t, err, "failed to get current working directory")
	require.NoError(t, os.Chdir(tempDir), "failed to change working directory")
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	// 5. Set the environment variable DIRECT_URL to your Supabase DSN.
	// Replace the below DSN with the one provided by your Supabase project.
	testDSN := "postgres://postgres.mwaowvtmopqrbknfojxs:Itish@Kiit@aws-0-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require"
	os.Setenv("DIRECT_URL", testDSN)

	// 6. Initialize the global DB by calling ConnectDB().
	err = ConnectDB()
	require.NoError(t, err, "ConnectDB failed")
	require.NotNil(t, DB, "DB should not be nil")

	// 7. For a fresh migration, drop the migrations table if it exists.
	err = DB.Exec("DROP TABLE IF EXISTS schema_migrations").Error
	require.NoError(t, err, "failed to drop schema_migrations table")

	if err != nil {
		t.Fatalf("Failed to drop migrations table: %v", err)
	}

	// 8. Run the migrations.
	err = Migrate()
	require.NoError(t, err, "Migrate() should not return an error")

	// 9. (Optional) Verify that the migration created the dummy table.
	var tableName string
	// PostgreSQL's to_regclass returns the table name if it exists.
	row := DB.Raw("SELECT to_regclass('public.dummy')").Row()
	err = row.Scan(&tableName)
	require.NoError(t, err, "failed to query for dummy table")
	require.Equal(t, "dummy", tableName, "expected dummy table to be created")
}
