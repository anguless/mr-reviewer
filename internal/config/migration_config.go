package config

import (
	"fmt"
	"os"
)

type MigrationConfig struct {
	MigrationsDir string
}

func loadMigrationConfig() (*MigrationConfig, error) {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		return nil, fmt.Errorf("MIGRATIONS_DIR is empty")
	}

	return &MigrationConfig{
		migrationsDir,
	}, nil
}
