package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port  int
	DbDsn string
}

func Load() (*Config, error) {
	config := &Config{
		Port: 8000,
	}

	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return nil, fmt.Errorf("Invalid PORT: %v", err)
		}
		config.Port = port
	}

	dbDsn := os.Getenv("DATABASE_DSN")
	if dbDsn == "" {
		return nil, fmt.Errorf("Specify DATABASE_DSN")
	}
	config.DbDsn = dbDsn

	return config, nil
}
