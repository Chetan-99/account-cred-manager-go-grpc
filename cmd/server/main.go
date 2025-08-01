package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/service"

	pb "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration : %v\n", err)
	}

	var accountServer pb.AccountServer

	if cfg.STORAGE_MODE == "MEM" {
		fmt.Printf("Starting service as MEM mode\n")
		accountServer = service.NewAccountsServiceMem()
	} else {
		fmt.Printf("Starting service as DB mode\n")
		var closeDB func() error
		accountServer, closeDB = service.NewAccountsServiceDB(cfg)
		defer closeDB()
	}

	g := NewGrpcServer(cfg)
	defer g.server.GracefulStop()

	pb.RegisterAccountServer(g.server, accountServer)

	go func() {
		g.Start()
	}()

	<-ctx.Done()
	fmt.Println("\nShutting down Application")
}
