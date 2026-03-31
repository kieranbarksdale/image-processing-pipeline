package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	ServerPort  string
	MinioEndpoint string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket string
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
	
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		panic("MINIO_ACCESS_KEY is not set")
	}
	
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		panic("MINIO_SECRET_KEY is not set")
	}

	MinioPort := os.Getenv("MINIO_PORT")
	if MinioPort == "" {
		panic("MINIO_PORT is not set")
	}
	
	MinioHost := os.Getenv("MINIO_HOST")
	if MinioHost == "" {
		panic("MINIO_HOST is not set")
	}

	MinioBucket := os.Getenv("MINIO_BUCKET")
	if MinioBucket == "" {
		panic("MINIO_BUCKET is not set")
	}

	minioEndpoint := fmt.Sprintf("%s:%s", MinioHost, MinioPort)

	return &Config{
		DatabaseURL:      databaseURL,
		RedisURL:         redisURL,
		ServerPort:       serverPort,
		MinioEndpoint:    minioEndpoint,
		MinioAccessKey:   minioAccessKey,
		MinioSecretKey:   minioSecretKey,
		MinioBucket:      MinioBucket,
	}
}
