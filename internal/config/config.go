package config

import "os"

type AppConfig struct {
	GRPC_PORT string
}

func LoadConfig() (*AppConfig, error) {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	return &AppConfig{
		GRPC_PORT: grpcPort,
	}, nil
}
