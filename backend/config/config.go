package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	OAuth     OAuthConfig
	CORS      CORSConfig
	Log       LogConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string
	Host string
	Env  string
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      time.Duration
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string
	Format string
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

// OAuthConfig holds OAuth provider settings.
type OAuthConfig struct {
	Google OAuthProviderConfig
	Naver  OAuthProviderConfig
	Kakao  OAuthProviderConfig
}

// OAuthProviderConfig holds individual OAuth provider credentials.
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// DSN returns the PostgreSQL connection string.
func (d *DatabaseConfig) DSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode +
		" TimeZone=Asia/Seoul"
}

// IsDevelopment returns true when running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true when running in production mode.
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// Load reads configuration from environment variables.
// It always attempts to load .env file (environment variables take precedence).
func Load() *Config {
	env := getEnv("ENV", "development")
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file, using environment variables", "error", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Env:  env,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "subkeep_user"),
			Password:        getEnv("DB_PASSWORD", "subkeep_dev_password"),
			Name:            getEnv("DB_NAME", "subkeep_db"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			TTL:      getEnvDuration("REDIS_TTL", 1*time.Hour),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", ""),
			Expiration:        getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshExpiration: getEnvDuration("JWT_REFRESH_EXPIRATION", 168*time.Hour),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			AllowedMethods: getEnvSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			Burst:             getEnvInt("RATE_LIMIT_BURST", 10),
		},
		OAuth: OAuthConfig{
			Google: OAuthProviderConfig{
				ClientID:     getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_GOOGLE_REDIRECT_URL", ""),
			},
			Naver: OAuthProviderConfig{
				ClientID:     getEnv("OAUTH_NAVER_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_NAVER_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_NAVER_REDIRECT_URL", ""),
			},
			Kakao: OAuthProviderConfig{
				ClientID:     getEnv("OAUTH_KAKAO_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_KAKAO_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_KAKAO_REDIRECT_URL", ""),
			},
		},
	}

	if cfg.JWT.Secret == "" && !cfg.IsDevelopment() {
		slog.Error("JWT_SECRET must be set in non-development environments")
		os.Exit(1)
	}

	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "subkeep-dev-jwt-secret-change-in-production"
		slog.Warn("using default JWT secret for development")
	}

	return cfg
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt returns an integer environment variable or a default value.
func getEnvInt(key string, defaultVal int) int {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		slog.Warn("invalid integer env var, using default", "key", key, "value", val, "default", defaultVal)
		return defaultVal
	}
	return n
}

// getEnvDuration returns a time.Duration environment variable or a default.
// Accepts formats: "24h", "30m", "5" (interpreted as minutes for DB_CONN_MAX_LIFETIME).
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}

	// Try parsing as a duration string first (e.g., "24h", "30m").
	d, err := time.ParseDuration(val)
	if err == nil {
		return d
	}

	// Fall back to interpreting as integer minutes.
	n, err := strconv.Atoi(val)
	if err != nil {
		slog.Warn("invalid duration env var, using default", "key", key, "value", val, "default", defaultVal)
		return defaultVal
	}
	return time.Duration(n) * time.Minute
}

// getEnvSlice returns a comma-separated environment variable as a string slice.
func getEnvSlice(key string, defaultVal []string) []string {
	val := getEnv(key, "")
	if val == "" {
		return defaultVal
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return defaultVal
	}
	return result
}
