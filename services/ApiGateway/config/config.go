package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
)

type Config struct {
	// Server
	AppPort string
	AppEnv  string

	// JWT
	JWTSecret string

	// CORS
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Service URLs
	UserServiceURL    string
	ProductServiceURL string
	CartServiceURL    string
	OrderServiceURL   string

	// Timeouts
	RequestTimeout time.Duration
	IdleTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	// Service name
	ServiceName string
}

func Load() (*Config, error) {
	// Try multiple paths for .env file
	envPaths := []string{
		filepath.Join("services/ApiGateway/config/.env"),
		filepath.Join("config/.env"),
		filepath.Join("./.env"),
	}

	var err error
	for _, envPath := range envPaths {
		err = godotenv.Load(envPath)
		if err == nil {
			logger.Infof("loaded .env file from: %s", envPath)
			break
		}
	}

	if err != nil {
		logger.Warnf("could not load .env file from any path: %v", err)
	}

	cfg := &Config{
		// Server
		AppPort: GetEnv("APP_PORT", "8080"),
		AppEnv:  GetEnv("APP_ENV", "development"),

		// JWT
		JWTSecret: GetEnv("JWT_SECRET", "your-secret-key-change-in-production"),

		// CORS
		AllowedOrigins: getEnvArray("ALLOWED_ORIGINS", []string{"*"}),
		AllowedMethods: getEnvArray("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		AllowedHeaders: getEnvArray("ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"}),

		// Rate Limiting
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   time.Duration(getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,

		// Service URLs
		UserServiceURL:    GetEnv("USER_SERVICE_URL", "localhost:50051"),
		ProductServiceURL: GetEnv("PRODUCT_SERVICE_URL", "localhost:50052"),
		CartServiceURL:    GetEnv("CART_SERVICE_URL", "localhost:50053"),
		OrderServiceURL:   GetEnv("ORDER_SERVICE_URL", "localhost:50054"),

		// Timeouts
		RequestTimeout: time.Duration(getEnvInt("REQUEST_TIMEOUT_SECONDS", 30)) * time.Second,
		IdleTimeout:    time.Duration(getEnvInt("IDLE_TIMEOUT_SECONDS", 120)) * time.Second,
		ReadTimeout:    time.Duration(getEnvInt("READ_TIMEOUT_SECONDS", 15)) * time.Second,
		WriteTimeout:   time.Duration(getEnvInt("WRITE_TIMEOUT_SECONDS", 15)) * time.Second,

		// Service
		ServiceName: GetEnv("SERVICE_NAME", "api-gateway"),
	}

	return cfg, nil
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value == "true" || value == "1" || value == "yes"
}

func getEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []string
	current := ""
	for _, char := range value {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}
