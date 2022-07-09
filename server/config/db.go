package config

import (
	"fmt"
)

func loadDBConfig(config *Config) {
	config.DBUser = getEnv("DB_USER")
	config.DBPassword = getEnv("DB_PASS")
	config.DBName = getEnv("DB_NAME")
	config.DBHost = getEnv("DB_HOST")
	config.DBPort = getEnv("DB_PORT")
}

func GetDBDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBPassword)
}
