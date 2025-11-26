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
	dbUsername := os.Getenv("POSTGRES_USER")
	if dbUsername == "" {
		return nil, fmt.Errorf("POSTGRES_USER is empty")
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is empty")
	}

	dBName := os.Getenv("POSTGRES_DB")
	if dBName == "" {
		return nil, fmt.Errorf("POSTGRES_DB is empty")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("POSTGRES_HOST is empty")
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("POSTGRES_PORT is empty")
	}

	return &DbConfig{
		Username:  dbUsername,
		Password:  dbPassword,
		TableName: dBName,
		Host:      dbHost,
		Port:      dbPort,
	}, nil
}
