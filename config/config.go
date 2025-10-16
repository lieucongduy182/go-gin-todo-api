package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBSSLMode          string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	expirationHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expirationHours = 24
	}

	AppConfig = &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "password"),
		DBName:             getEnv("DB_NAME", "todo_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "default-secret"),
		JWTExpirationHours: expirationHours,
	}

	log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if val, exist := os.LookupEnv(key); exist {
		return val
	}

	return defaultValue
}
