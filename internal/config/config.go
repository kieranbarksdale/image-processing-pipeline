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

	getEnvOrPanic := func(key string) string {
        val := os.Getenv(key)
        if val == "" {
            panic(fmt.Sprintf("CRITICAL: Environment variable %s is not set", key))
        }
        return val
    }
	
	migrationsPath := getEnvOrPanic("MIGRATIONS_PATH")
	redisHost := getEnvOrPanic("REDIS_HOST")
	redisPort := getEnvOrPanic("REDIS_PORT")
	serverPort := ":" + getEnvOrPanic("API_PORT")
	minioAccessKey := getEnvOrPanic("MINIO_ACCESS_KEY")
	minioSecretKey := getEnvOrPanic("MINIO_SECRET_KEY")
	MinioPort := getEnvOrPanic("MINIO_PORT")
	MinioHost := getEnvOrPanic("MINIO_HOST")
	MinioBucket := getEnvOrPanic("MINIO_BUCKET")

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
