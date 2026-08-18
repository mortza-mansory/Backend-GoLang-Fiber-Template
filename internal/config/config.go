// Package config loads and validates application configuration from the
// environment. Configuration is intentionally explicit and boring: every
// setting maps to an environment variable (see .env.example) and falls back
// to a sane default for local development.
//
// Configuration is loaded once at startup and passed to the rest of the
// application through constructors. It is not a global singleton.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config is the root application configuration.
type Config struct {
	App       App
	Server    Server
	Database  Database
	Redis     Redis
	Logger    Logger
	OTel      OTel
	External  External
	Security  Security
	RateLimit RateLimit
}

// App holds application-level metadata.
type App struct {
	Name        string
	Environment string
	Version     string
	Debug       bool
}

// Server holds HTTP server settings.
type Server struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	BodyLimit       int
}

// Database holds PostgreSQL connection settings.
type Database struct {
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	SSLMode          string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	ConnMaxIdleTime  time.Duration
	RunMigrations    bool
	MigrationsSource string
}

// Redis holds Redis connection settings.
type Redis struct {
	Host         string
	Port         string
	Password     string
	DB           int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Logger holds structured logging settings.
type Logger struct {
	Level string
	JSON  bool
}

// OTel holds OpenTelemetry settings.
type OTel struct {
	Enabled         bool
	ServiceName     string
	ServiceVersion  string
	Environment     string
	OTLPEndpoint    string
	TracesEndpoint  string
	MetricsEndpoint string
	Insecure        bool
	SamplingRatio   float64
	BatchTimeout    time.Duration
}

// External groups generic external service settings.
type External struct {
	Email   EmailService
	Storage StorageService
}

// EmailService holds placeholder SMTP settings.
type EmailService struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// StorageService holds placeholder object storage settings.
type StorageService struct {
	Provider string
	Bucket   string
	Region   string
}

// Security holds token/signing settings.
type Security struct {
	JWTSecret     string
	JWTExpiry     time.Duration
	JWTIssuer     string
	JWTAlgorithm  string
	BcryptCost    int
	AllowedOrigin string
}

// RateLimit holds global rate limiting settings.
type RateLimit struct {
	Enabled    bool
	Max        int
	Expiration time.Duration
}

// Load reads configuration from the process environment.
func Load() (*Config, error) {
	cfg := &Config{
		App: App{
			Name:        getEnv("APP_NAME", "go-fiber-template"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "0.1.0"),
			Debug:       getEnvBool("APP_DEBUG", false),
		},
		Server: Server{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
			BodyLimit:       getEnvInt("SERVER_BODY_LIMIT", 4*1024*1024),
		},
		Database: Database{
			Host:             getEnv("DB_HOST", "localhost"),
			Port:             getEnv("DB_PORT", "5432"),
			User:             getEnv("DB_USER", "postgres"),
			Password:         getEnv("DB_PASSWORD", "postgres"),
			Name:             getEnv("DB_NAME", "app"),
			SSLMode:          getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:     getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:     getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime:  getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime:  getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
			RunMigrations:    getEnvBool("DB_RUN_MIGRATIONS", true),
			MigrationsSource: getEnv("DB_MIGRATIONS_SOURCE", "file://migrations"),
		},
		Redis: Redis{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvInt("REDIS_DB", 0),
			MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		Logger: Logger{
			Level: getEnv("LOG_LEVEL", "info"),
			JSON:  getEnvBool("LOG_JSON", false),
		},
		OTel: OTel{
			Enabled:         getEnvBool("OTEL_ENABLED", false),
			ServiceName:     getEnv("OTEL_SERVICE_NAME", "go-fiber-template"),
			ServiceVersion:  getEnv("OTEL_SERVICE_VERSION", "0.1.0"),
			Environment:     getEnv("OTEL_ENVIRONMENT", "development"),
			OTLPEndpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
			TracesEndpoint:  getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", ""),
			MetricsEndpoint: getEnv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", ""),
			Insecure:        getEnvBool("OTEL_INSECURE", true),
			SamplingRatio:   getEnvFloat("OTEL_SAMPLING_RATIO", 1.0),
			BatchTimeout:    getEnvDuration("OTEL_BATCH_TIMEOUT", 5*time.Second),
		},
		External: External{
			Email: EmailService{
				Host:     getEnv("SMTP_HOST", ""),
				Port:     getEnvInt("SMTP_PORT", 587),
				User:     getEnv("SMTP_USER", ""),
				Password: getEnv("SMTP_PASSWORD", ""),
				From:     getEnv("SMTP_FROM", ""),
			},
			Storage: StorageService{
				Provider: getEnv("STORAGE_PROVIDER", "local"),
				Bucket:   getEnv("STORAGE_BUCKET", ""),
				Region:   getEnv("STORAGE_REGION", ""),
			},
		},
		Security: Security{
			JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpiry:     getEnvDuration("JWT_EXPIRY", 15*time.Minute),
			JWTIssuer:     getEnv("JWT_ISSUER", "go-fiber-template"),
			JWTAlgorithm:  getEnv("JWT_ALGORITHM", "HS256"),
			BcryptCost:    getEnvInt("BCRYPT_COST", 10),
			AllowedOrigin: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		RateLimit: RateLimit{
			Enabled:    getEnvBool("RATE_LIMIT_ENABLED", true),
			Max:        getEnvInt("RATE_LIMIT_MAX", 100),
			Expiration: getEnvDuration("RATE_LIMIT_EXPIRATION", 1*time.Minute),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// validate performs basic validation and returns human-readable errors.
func (c *Config) validate() error {
	if c.Security.JWTSecret == "" || c.Security.JWTSecret == "change-me-in-production" {
		// Only warn in non-production to avoid a hard failure during local dev.
		if c.App.Environment == "production" {
			return fmt.Errorf("config: JWT_SECRET must be set in production")
		}
	}
	return nil
}

// Address returns the host:port the HTTP server should bind to.
func (s Server) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// DSN builds a PostgreSQL connection string.
func (d Database) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Addr returns the Redis host:port.
func (r Redis) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(v), "%d", &i); err != nil {
		return fallback
	}
	return i
}

func getEnvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	f := 0.0
	if _, err := fmt.Sscanf(strings.TrimSpace(v), "%f", &f); err != nil {
		return fallback
	}
	return f
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	d, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return fallback
	}
	return d
}
