package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig        *DbConfig
	MigrationConfig *MigrationConfig
	AppConfig       *AppConfig
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	dbConfig, err := loadDbConfig()
	if err != nil {
		return nil, err
	}

	migrationConfig, err := loadMigrationConfig()
	if err != nil {
		return nil, err
	}

	appConfig := loadAppConfig()

	return &Config{
		DbConfig:        dbConfig,
		MigrationConfig: migrationConfig,
		AppConfig:       appConfig,
	}, nil
}
