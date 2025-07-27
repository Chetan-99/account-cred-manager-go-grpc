package main

import (
	"log"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/service"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/store"

	pb "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration : %v", err)
	}

	g := NewGrpcServer(cfg)

	accounts_store := store.NewAccountStore()

	accountService := service.NewAccountsService(accounts_store)
	pb.RegisterAccountServer(g.server, accountService)

	g.Start()
}
