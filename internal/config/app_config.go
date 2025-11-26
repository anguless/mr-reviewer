package config

import (
	"os"
)

type AppConfig struct {
	AppPort string
}

var defaultAppPort = "8080"

func loadAppConfig() (*AppConfig, error) {
	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		appPort = defaultAppPort
	}

	return &AppConfig{
		AppPort: appPort,
	}, nil
}
