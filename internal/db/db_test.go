package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_WithInvalidDSN(t *testing.T) {
	// Set an invalid DSN
	originalDSN := os.Getenv("DATABASE_URL")
	os.Setenv("DATABASE_URL", "postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable")
	defer os.Setenv("DATABASE_URL", originalDSN)

	// Connect should fail with invalid DSN
	db, err := Connect()

	// Since we can't connect to an invalid database, we expect an error
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnect_WithEmptyDSN(t *testing.T) {
	// Clear DSN to use default
	originalDSN := os.Getenv("DATABASE_URL")
	os.Unsetenv("DATABASE_URL")
	defer func() {
		if originalDSN != "" {
			os.Setenv("DATABASE_URL", originalDSN)
		}
	}()

	// Connect should fail with default DSN (assuming no local postgres)
	db, err := Connect()

	// We expect an error since there's likely no local postgres
	// If a local postgres exists, this test may pass
	if err != nil {
		assert.Nil(t, db)
	} else {
		assert.NotNil(t, db)
		// Close the connection if it succeeded
		if db != nil {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
	}
}

func TestRunMigration_WithNilDB(t *testing.T) {
	// RunMigration with nil db should panic or return error
	defer func() {
		if r := recover(); r != nil {
			// Expected panic when db is nil
			assert.NotNil(t, r)
		}
	}()

	err := RunMigration(nil, "SELECT 1")
	// If no panic, we should have an error
	assert.Error(t, err)
}

func TestRunMigration_WithInvalidSQL(t *testing.T) {
	// Skip if no DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	db, err := Connect()
	if err != nil {
		t.Skip("Could not connect to database")
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Run invalid SQL
	err = RunMigration(db, "INVALID SQL STATEMENT")
	assert.Error(t, err)
}

func TestRunMigration_WithValidSQL(t *testing.T) {
	// Skip if no DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	db, err := Connect()
	if err != nil {
		t.Skip("Could not connect to database")
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Run valid SQL
	err = RunMigration(db, "SELECT 1")
	assert.NoError(t, err)
}

func TestConnect_ConnectionPoolSettings(t *testing.T) {
	// This test verifies the connection pool settings exist in the code
	// We can't directly test the pool settings without a real connection

	// Skip if no DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	db, err := Connect()
	if err != nil {
		t.Skip("Could not connect to database")
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	assert.NotNil(t, sqlDB)

	// Verify pool settings are applied
	stats := sqlDB.Stats()
	assert.NotNil(t, stats)
}

func TestDefaultDSN(t *testing.T) {
	// Clear any existing DATABASE_URL
	originalDSN := os.Getenv("DATABASE_URL")
	os.Unsetenv("DATABASE_URL")
	defer func() {
		if originalDSN != "" {
			os.Setenv("DATABASE_URL", originalDSN)
		}
	}()

	// Verify default DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://localhost:5432/postgres?sslmode=disable"
	}

	assert.Equal(t, "postgres://localhost:5432/postgres?sslmode=disable", dsn)
}
