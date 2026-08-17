package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa todas as configurações da aplicação.
type Config struct {
	App AppConfig
	DB  DatabaseConfig
	JWT JWTConfig
}

type AppConfig struct {
	Name       string
	Env        string
	Port       string
	BcryptCost int
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// func Load() *Config {
// 	_ = godotenv.Load()

// 	return &Config{
// 		App: AppConfig{
// 			Name:       getEnv("APP_NAME", "go-api-base"),
// 			Env:        getEnv("APP_ENV", "development"),
// 			Port:       getEnv("APP_PORT", "8080"),
// 			BcryptCost: getEnvAsInt("BCRYPT_COST", 10),
// 		},
// 		DB: DatabaseConfig{
// 			Host:         getEnv("DB_HOST", "localhost"),
// 			Port:         getEnv("DB_PORT", "5432"),
// 			User:         getEnv("DB_USER", "postgres"),
// 			Password:     getEnv("DB_PASSWORD", "postgres"),
// 			Name:         getEnv("DB_NAME", "go_api_base"),
// 			SSLMode:      getEnv("DB_SSLMODE", "disable"),
// 			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
// 			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
// 		},
// 		JWT: JWTConfig{
// 			Secret:          getEnv("JWT_SECRET", "change-me"),
// 			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
// 			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TTL", 168*time.Hour),
// 		},
// 	}
// }

// func Load() *Config {
// 	_ = godotenv.Load()

// 	return &Config{
// 		App: AppConfig{
// 			Name:       getEnv("APP_NAME"),
// 			Env:        getEnv("APP_ENV"),
// 			Port:       getEnv("APP_PORT"),
// 			BcryptCost: getEnvAsInt("BCRYPT_COST"),
// 		},
// 		DB: DatabaseConfig{
// 			Host:         getEnv("DB_HOST"),
// 			Port:         getEnv("DB_PORT"),
// 			User:         getEnv("DB_USER"),
// 			Password:     getEnv("DB_PASSWORD"),
// 			Name:         getEnv("DB_NAME"),
// 			SSLMode:      getEnv("DB_SSLMODE"),
// 			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS"),
// 			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS"),
// 		},
// 		JWT: JWTConfig{
// 			Secret:          getEnv("JWT_SECRET"),
// 			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TTL"),
// 			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TTL"),
// 		},
// 	}
// }

// func (d DatabaseConfig) DSN() string {
// 	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port +
// 		"/" + d.Name + "?sslmode=" + d.SSLMode
// }

// func getEnv(key, fallback string) string {
// 	if v, ok := os.LookupEnv(key); ok && v != "" {
// 		return v
// 	}
// 	return fallback
// }

// func getEnvAsInt(key string, fallback int) int {
// 	if v, ok := os.LookupEnv(key); ok {
// 		if i, err := strconv.Atoi(v); err == nil {
// 			return i
// 		}
// 	}
// 	return fallback
// }

// func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
// 	if v, ok := os.LookupEnv(key); ok {
// 		if d, err := time.ParseDuration(v); err == nil {
// 			return d
// 		}
// 	}
// 	return fallback
// }

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:       getEnv("APP_NAME"),
			Env:        getEnv("APP_ENV"),
			Port:       getEnv("APP_PORT"),
			BcryptCost: getEnvAsInt("BCRYPT_COST"),
		},
		DB: DatabaseConfig{
			Host:         getEnv("DB_HOST"),
			Port:         getEnv("DB_PORT"),
			User:         getEnv("DB_USER"),
			Password:     getEnv("DB_PASSWORD"),
			Name:         getEnv("DB_NAME"),
			SSLMode:      getEnv("DB_SSLMODE"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET"),
			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TTL"),
			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TTL"),
		},
	}
}

func (d DatabaseConfig) DSN() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port +
		"/" + d.Name + "?sslmode=" + d.SSLMode
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok || value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}

	return value
}

func getEnvAsInt(key string) int {
	value := getEnv(key)

	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("environment variable %s must be an integer: %v", key, err)
	}

	return result
}

func getEnvAsDuration(key string) time.Duration {
	value := getEnv(key)

	result, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("environment variable %s must be a valid duration: %v", key, err)
	}

	return result
}
