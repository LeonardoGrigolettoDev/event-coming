package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App       AppConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	WhatsApp  WhatsAppConfig
	OSRM      OSRMConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig holds PostgreSQL connection configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxConnAge      time.Duration `mapstructure:"max_conn_age"`
	PoolTimeout     time.Duration `mapstructure:"pool_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	AccessSecret       string        `mapstructure:"access_secret"`
	RefreshSecret      string        `mapstructure:"refresh_secret"`
	AccessTokenTTL     time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL    time.Duration `mapstructure:"refresh_token_ttl"`
	Issuer             string        `mapstructure:"issuer"`
}

// WhatsAppConfig holds WhatsApp Cloud API configuration
type WhatsAppConfig struct {
	VerifyToken    string `mapstructure:"verify_token"`
	AppSecret      string `mapstructure:"app_secret"`
	AccessToken    string `mapstructure:"access_token"`
	PhoneNumberID  string `mapstructure:"phone_number_id"`
	BusinessID     string `mapstructure:"business_id"`
	APIVersion     string `mapstructure:"api_version"`
	BaseURL        string `mapstructure:"base_url"`
}

// OSRMConfig holds OSRM routing service configuration
type OSRMConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// Load reads configuration from environment variables and files
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from environment
	v.AutomaticEnv()
	v.SetEnvPrefix("EVENT_COMING")

	// Read from .env file if exists
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	_ = v.ReadInConfig() // Ignore error if .env doesn't exist

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "event-coming")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.database", "event_coming")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)
	v.SetDefault("database.max_conn_lifetime", 1*time.Hour)
	v.SetDefault("database.max_conn_idle_time", 30*time.Minute)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.max_conn_age", 0)
	v.SetDefault("redis.pool_timeout", 4*time.Second)
	v.SetDefault("redis.idle_timeout", 5*time.Minute)

	// JWT defaults
	v.SetDefault("jwt.access_secret", "change-me-in-production")
	v.SetDefault("jwt.refresh_secret", "change-me-in-production")
	v.SetDefault("jwt.access_token_ttl", 15*time.Minute)
	v.SetDefault("jwt.refresh_token_ttl", 7*24*time.Hour)
	v.SetDefault("jwt.issuer", "event-coming")

	// WhatsApp defaults
	v.SetDefault("whatsapp.verify_token", "")
	v.SetDefault("whatsapp.app_secret", "")
	v.SetDefault("whatsapp.access_token", "")
	v.SetDefault("whatsapp.phone_number_id", "")
	v.SetDefault("whatsapp.business_id", "")
	v.SetDefault("whatsapp.api_version", "v18.0")
	v.SetDefault("whatsapp.base_url", "https://graph.facebook.com")

	// OSRM defaults
	v.SetDefault("osrm.enabled", false)
	v.SetDefault("osrm.base_url", "http://localhost:5000")
	v.SetDefault("osrm.timeout", 10*time.Second)
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
