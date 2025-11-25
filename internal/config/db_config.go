package config

import (
	"fmt"
	"os"
)

type DbConfig struct {
	Username  string
	Password  string
	TableName string
	Host      string
	Port      string
}

func loadDbConfig() (*DbConfig, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		return nil, fmt.Errorf("DB_USERNAME is empty")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is empty")
	}

	dBName := os.Getenv("DB_NAME")
	if dBName == "" {
		return nil, fmt.Errorf("DB_NAME is empty")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST is empty")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("DB_PORT is empty")
	}

	return &DbConfig{
		Username: dbUsername,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
	}, nil
}
