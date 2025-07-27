package service

import (
	"context"
	"fmt"

	pb "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountsService struct {
	pb.UnimplementedAccountServer
	accounts_store *store.AccountsStore
}

func NewAccountsService(s *store.AccountsStore) *AccountsService {
	return &AccountsService{
		accounts_store: s,
	}
}

func (p *AccountsService) CreateAccount(ctx context.Context, in *pb.AccountInputRequest) (*pb.CreateAccountResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Creating account with account id - %d\n", account_id)
	res, err := p.accounts_store.CreateAccount(account_id)
	if err != nil {
		fmt.Printf("Failed to create account with account id - %d - error - %v", account_id, err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.CreateAccountResponse{AccountId: account_id, SessionToken: res}, nil
}

func (p *AccountsService) GetToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Getting token for account_id = %d\n", account_id)
	res, err := p.accounts_store.GetToken(account_id)
	if err != nil {
		fmt.Printf("Failed to get token for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.TokenResponse{SessionToken: res}, nil
}

func (p *AccountsService) IsTokenExpired(ctx context.Context, in *pb.AccountInputRequest) (*pb.IsTokenExpiredResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Checking if token expired for account_id = %d\n", account_id)
	res, err := p.accounts_store.IsTokenExpired(account_id)
	if err != nil {
		fmt.Printf("Failed to check token expiry for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.Canceled, err.Error())
	}
	return &pb.IsTokenExpiredResponse{Expired: res}, nil
}

func (p *AccountsService) RegenerateToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Re-generating token for account_id = %d\n", account_id)
	res, err := p.accounts_store.RegenerateToken(account_id)
	if err != nil {
		fmt.Printf("Failed to re-generate token for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.Canceled, err.Error())
	}
	return &pb.TokenResponse{SessionToken: res}, nil
}
