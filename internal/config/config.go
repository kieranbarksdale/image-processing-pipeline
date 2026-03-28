package config

import "os"

type Config struct {
	DatabaseURL string
	RedisURL    string
	ServerPort  string
	MinioURL    string
	MinioAccessKey string
	MinioSecretKey string
}

func Load() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL is not set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL is not set")
	}

	serverPort := os.Getenv("API_PORT")
	if serverPort == "" {
		panic("API_PORT is not set")
	}
	serverPort = ":" + serverPort
	
	minioURL := os.Getenv("MINIO_URL")
	if minioURL == "" {
		panic("MINIO_URL is not set")
	}
	
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		panic("MINIO_ACCESS_KEY is not set")
	}
	
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		panic("MINIO_SECRET_KEY is not set")
	}

	return &Config{
		DatabaseURL:      databaseURL,
		RedisURL:         redisURL,
		ServerPort:       serverPort,
		MinioURL:         minioURL,
		MinioAccessKey:   minioAccessKey,
		MinioSecretKey:   minioSecretKey,
	}
}
