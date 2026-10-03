package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FirebaseProjectID      string
	InstanceConnectionName string
	DatabaseName           string
	DatabaseUser           string
	DatabasePassword       string
	HTTPAddr               string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("cargar .env: %w", err)
	}

	config := Config{
		FirebaseProjectID:      os.Getenv("FIREBASE_PROJECT_ID"),
		InstanceConnectionName: os.Getenv("INSTANCE_CONNECTION_NAME"),
		DatabaseName:           os.Getenv("DB_NAME"),
		DatabaseUser:           os.Getenv("DB_USER"),
		DatabasePassword:       os.Getenv("DB_PASSWORD"),
		HTTPAddr:               os.Getenv("HTTP_ADDR"),
	}
	if config.FirebaseProjectID == "" ||
		config.InstanceConnectionName == "" ||
		config.DatabaseName == "" ||
		config.DatabaseUser == "" ||
		config.DatabasePassword == "" {
		return Config{}, errors.New("FIREBASE_PROJECT_ID, INSTANCE_CONNECTION_NAME, DB_NAME, DB_USER y DB_PASSWORD son obligatorios")
	}
	if os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") != "" && os.Getenv("APP_ENV") != "development" {
		return Config{}, errors.New("FIREBASE_AUTH_EMULATOR_HOST solo se permite con APP_ENV=development")
	}
	if config.HTTPAddr == "" {
		config.HTTPAddr = "127.0.0.1:8080"
	}
	return config, nil
}
