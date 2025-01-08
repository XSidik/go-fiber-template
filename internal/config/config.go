package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  int
	RedisDB    int
	JWTSecret  string
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %v", name, err)
	}
	return value
}

func GetConfig() Config {
	// Load .env file debug
	// err := godotenv.Load("../.env")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "dbname"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnvAsInt("REDIS_PORT", 6379),
		RedisDB:    getEnvAsInt("REDIS_DB", 0),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}
}
