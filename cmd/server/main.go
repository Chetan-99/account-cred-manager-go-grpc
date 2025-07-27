package main

import (
	"fmt"
	"log"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/service"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/store"

	pb "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration : %v\n", err)
	}

	var accountServer pb.AccountServer
	if cfg.STORAGE_MODE == "MEM" {
		fmt.Printf("Starting service as MEM mode\n")
		mem_store := store.NewAccountStore()
		accountServer = service.NewAccountsServiceMem(mem_store)
	} else {
		fmt.Printf("Starting service as DB mode\n")
		db := store.NewBadgerDB(cfg)
		defer db.Close()
		accountServer = service.NewAccountsServiceDB(db)
	}

	g := NewGrpcServer(cfg)

	pb.RegisterAccountServer(g.server, accountServer)

	g.Start()
}
