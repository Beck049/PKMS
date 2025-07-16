package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	SearchPath string
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "pkms"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		SearchPath: getEnv("SEARCH_PATH", "/app/articles"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
