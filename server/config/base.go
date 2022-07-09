package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppEnv string

const (
	AppEnvDevelopment = "development"
	AppEnvProdution   = "production"
)

type Config struct {
	DBUser       string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       string
	AppPort      string
	AppEnv       AppEnv
	AppURL       string
	JwtSecretKey string
}

var config = Config{}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	loadDBConfig(&config)
	loadAppConfig(&config)
	loadJWTConfig(&config)
}

func GetConfig() Config {
	return config
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		log.Fatalf("Missing or invalid environment key: '%s'", key)
	}

	return value
}
func getEnvWithDefault(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
