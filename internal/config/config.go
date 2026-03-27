package config

import "os"

type Config struct {
	DatabaseURL string
	RedisURL    string
	ServerPort  string
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

	return &Config{
		DatabaseURL: databaseURL,
		RedisURL:    redisURL,
		ServerPort:  serverPort,
	}
}
