package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Clear any existing environment variables
	oldEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range oldEnv {
			kv := splitEnv(e)
			if len(kv) == 2 {
				os.Setenv(kv[0], kv[1])
			}
		}
	}()

	t.Run("loads default configuration", func(t *testing.T) {
		cfg, err := Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// Check defaults
		assert.Equal(t, "event-coming", cfg.App.Name)
		assert.Equal(t, "development", cfg.App.Environment)
		assert.True(t, cfg.App.Debug)

		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)

		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "postgres", cfg.Database.User)

		assert.Equal(t, "localhost", cfg.Redis.Host)
		assert.Equal(t, 6379, cfg.Redis.Port)

		assert.Equal(t, "event-coming", cfg.JWT.Issuer)
	})

	t.Run("loads from environment variables", func(t *testing.T) {
		os.Setenv("EVENT_COMING_SERVER_PORT", "9090")
		os.Setenv("EVENT_COMING_DATABASE_HOST", "db.example.com")
		os.Setenv("EVENT_COMING_REDIS_HOST", "redis.example.com")
		defer func() {
			os.Unsetenv("EVENT_COMING_SERVER_PORT")
			os.Unsetenv("EVENT_COMING_DATABASE_HOST")
			os.Unsetenv("EVENT_COMING_REDIS_HOST")
		}()

		cfg, err := Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, 9090, cfg.Server.Port)
		assert.Equal(t, "db.example.com", cfg.Database.Host)
		assert.Equal(t, "redis.example.com", cfg.Redis.Host)
	})
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "standard configuration",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secret",
				Database: "mydb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable",
		},
		{
			name: "production configuration",
			config: DatabaseConfig{
				Host:     "db.production.com",
				Port:     5433,
				User:     "admin",
				Password: "prod-password",
				Database: "event_coming",
				SSLMode:  "require",
			},
			expected: "host=db.production.com port=5433 user=admin password=prod-password dbname=event_coming sslmode=require",
		},
		{
			name: "empty password",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "",
				Database: "test",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password= dbname=test sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetDSN()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedisConfig_GetRedisAddr(t *testing.T) {
	tests := []struct {
		name     string
		config   RedisConfig
		expected string
	}{
		{
			name: "localhost default",
			config: RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
			expected: "localhost:6379",
		},
		{
			name: "production redis",
			config: RedisConfig{
				Host: "redis.production.com",
				Port: 6380,
			},
			expected: "redis.production.com:6380",
		},
		{
			name: "custom port",
			config: RedisConfig{
				Host: "127.0.0.1",
				Port: 16379,
			},
			expected: "127.0.0.1:16379",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetRedisAddr()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Defaults(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)

	// App defaults
	assert.Equal(t, "event-coming", cfg.App.Name)
	assert.Equal(t, "development", cfg.App.Environment)
	assert.True(t, cfg.App.Debug)

	// Server defaults
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, 60*time.Second, cfg.Server.IdleTimeout)

	// Database defaults
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "postgres", cfg.Database.Password)
	assert.Equal(t, "event_coming", cfg.Database.Database)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, int32(25), cfg.Database.MaxConns)
	assert.Equal(t, int32(5), cfg.Database.MinConns)

	// Redis defaults
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, 10, cfg.Redis.PoolSize)
	assert.Equal(t, 5, cfg.Redis.MinIdleConns)

	// JWT defaults
	assert.Equal(t, "change-me-in-production", cfg.JWT.AccessSecret)
	assert.Equal(t, "change-me-in-production", cfg.JWT.RefreshSecret)
	assert.Equal(t, "event-coming", cfg.JWT.Issuer)
	assert.Equal(t, 15*time.Minute, cfg.JWT.AccessExpiresIn)
	assert.Equal(t, 7*24*time.Hour, cfg.JWT.RefreshExpiresIn)

	// WhatsApp defaults
	assert.Equal(t, "v18.0", cfg.WhatsApp.APIVersion)
	assert.Equal(t, "https://graph.facebook.com", cfg.WhatsApp.BaseURL)
	assert.Equal(t, "event-coming-webhook-token", cfg.WhatsApp.WebhookVerifyToken)

	// OSRM defaults
	assert.False(t, cfg.OSRM.Enabled)
	assert.Equal(t, "http://localhost:5000", cfg.OSRM.BaseURL)
	assert.Equal(t, 10*time.Second, cfg.OSRM.Timeout)
}

// Helper function to split environment variable
func splitEnv(e string) []string {
	for i := 0; i < len(e); i++ {
		if e[i] == '=' {
			return []string{e[:i], e[i+1:]}
		}
	}
	return []string{e}
}
