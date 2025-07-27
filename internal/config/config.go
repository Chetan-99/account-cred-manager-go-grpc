package config

import (
	"os"
)

const (
	D_GRPC_PORT    = "50051"
	D_DB_PATH      = "./badger.db"
	D_STORAGE_MODE = "MEM"
)

type AppConfig struct {
	GRPC_PORT    string
	DB_PATH      string
	STORAGE_MODE string
}

func LoadConfig() (*AppConfig, error) {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = D_GRPC_PORT
	}

	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = D_DB_PATH
	}

	storage_path := os.Getenv("STORAGE_MODE")
	if storage_path == "" {
		storage_path = D_STORAGE_MODE
	}
	if storage_path != "MEM" && storage_path != "DB" {
		storage_path = D_STORAGE_MODE
	}
	// Ensure STORAGE_MODE is either "MEM" or "DB"
	if storage_path != "MEM" && storage_path != "DB" {
		storage_path = D_STORAGE_MODE
	}
	return &AppConfig{
		GRPC_PORT:    grpcPort,
		DB_PATH:      db_path,
		STORAGE_MODE: storage_path,
	}, nil
}
