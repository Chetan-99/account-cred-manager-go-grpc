package main

import (
	"log"
	"net"
	"os"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
}

func NewGrpcServer(cfg *config.AppConfig) *GRPCServer {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	s := grpc.NewServer()
	return &GRPCServer{
		server: s,
		port:   cfg.GRPC_PORT,
	}
}

func (s *GRPCServer) Start() {
	log.Printf("Starting grpc server on port - %s", s.port)
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("failed to start grpc server - error - %v", err)
	}

	log.Printf("grpc server listening at %v", lis.Addr())
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
