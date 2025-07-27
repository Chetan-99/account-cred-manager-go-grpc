package test

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	v1 "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddr           = "localhost:50051" // Change if your server runs elsewhere
	concurrentRequests = 500
)

func TestAccountEndpointsLoad(t *testing.T) {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := v1.NewAccountClient(conn)

	var wg sync.WaitGroup
	start := time.Now()

	accountIDs := make([]int32, concurrentRequests)
	for i := range accountIDs {
		accountIDs[i] = rand.Int31n(1_000_000)
	}

	// Create accounts concurrently
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			resp, err := client.CreateAccount(context.Background(), &v1.AccountInputRequest{AccountId: accountIDs[idx]})
			if err != nil {
				log.Printf("CreateAccount error for %d: %v", accountIDs[idx], err)
				return
			}
			// GetToken
			_, err = client.GetToken(context.Background(), &v1.AccountInputRequest{AccountId: resp.AccountId})
			if err != nil {
				log.Printf("GetToken error for %d: %v", resp.AccountId, err)
			}
			// RegenerateToken
			_, err = client.RegenerateToken(context.Background(), &v1.AccountInputRequest{AccountId: resp.AccountId})
			if err != nil {
				log.Printf("RegenerateToken error for %d: %v", resp.AccountId, err)
			}
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	t.Logf("Processed %d accounts in %v", concurrentRequests, duration)
}
