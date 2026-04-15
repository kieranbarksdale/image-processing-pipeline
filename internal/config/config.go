package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	MigrationsPath string
	RedisURL    string
	RedisHost   string
	RedisPort   string
	ServerPort  string
	MinioEndpoint string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket string
}

func Load() *Config {

	getEnvOrPanic := func(key string, fallback string) string {
        val := os.Getenv(key)
        if val == "" {
            return fallback
        }
        return val
    }
	
	migrationsPath := getEnvOrPanic("MIGRATIONS_PATH", "file:///app/internal/db/migrations")
	redisHost := getEnvOrPanic("REDIS_HOST", "redis")
	redisPort := getEnvOrPanic("REDIS_PORT", "6379")
	serverPort := ":" + getEnvOrPanic("API_PORT", "8080")
	minioAccessKey := getEnvOrPanic("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnvOrPanic("MINIO_SECRET_KEY", "minioadmin")
	MinioPort := getEnvOrPanic("MINIO_PORT", "9000")
	MinioHost := getEnvOrPanic("MINIO_HOST", "minio")
	MinioBucket := getEnvOrPanic("MINIO_BUCKET", "images")

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnvOrPanic("POSTGRES_USER"),
		getEnvOrPanic("POSTGRES_PASSWORD"),
		getEnvOrPanic("POSTGRES_HOST"),
		getEnvOrPanic("POSTGRES_PORT"),
		getEnvOrPanic("POSTGRES_DB"),
	)

	redisURL := fmt.Sprintf(
		"redis://%s:%s",
		redisHost,
		redisPort,
	)

	minioEndpoint := fmt.Sprintf("%s:%s", MinioHost, MinioPort)

	return &Config{
		DatabaseURL:      databaseURL,
		MigrationsPath:   migrationsPath,
		RedisURL:         redisURL,
		RedisHost:        redisHost,
		RedisPort:        redisPort,
		ServerPort:       serverPort,
		MinioEndpoint:    minioEndpoint,
		MinioAccessKey:   minioAccessKey,
		MinioSecretKey:   minioSecretKey,
		MinioBucket:      MinioBucket,
	}
}
