package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/chetan-99/account-cred-manager-go-grpc/api/proto/v1"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountsServiceDB struct {
	pb.UnimplementedAccountServer
	db *store.DbHandler
}

func NewAccountsServiceDB(cfg *config.AppConfig) (*AccountsServiceDB, func() error) {
	db := store.NewBadgerDB(cfg)
	return &AccountsServiceDB{
			db: db,
		}, func() error {
			return db.Close()
		}
}

func (p *AccountsServiceDB) CreateAccount(ctx context.Context, in *pb.AccountInputRequest) (*pb.CreateAccountResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Creating account with account id - %d\n", account_id)

	account := store.NewAccount(account_id)
	account_buf, err := account.Encode()
	if err != nil {
		fmt.Printf("failed to encode account id - %d - error - %v", account_id, err)
		return nil, err
	}

	err = p.db.Add_KV(account_id, account_buf)
	if err != nil {
		fmt.Printf("failed to create account - %d - error - %s", account_id, err)
		return nil, status.Error(codes.Internal, "Failed to create Account")
	}

	return &pb.CreateAccountResponse{AccountId: account_id, SessionToken: account.SessionToken}, nil
}

func (p *AccountsServiceDB) GetToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Getting token for account_id = %d\n", account_id)

	account, err := getAndDecodeAccount(account_id, p.db)
	if err != nil {
		fmt.Printf("Failed to get account - %v", err)
		return nil, status.Error(codes.Internal, "Failed to Get token")
	}

	return &pb.TokenResponse{SessionToken: account.GetToken()}, nil
}

func (p *AccountsServiceDB) IsTokenExpired(ctx context.Context, in *pb.AccountInputRequest) (*pb.IsTokenExpiredResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Checking if token expired for account_id = %d\n", account_id)

	account, err := getAndDecodeAccount(account_id, p.db)
	if err != nil {
		fmt.Printf("Failed to get token - %v", err)
		return nil, status.Error(codes.Internal, "Failed to Get token")
	}

	return &pb.IsTokenExpiredResponse{Expired: account.IsTokenExpired()}, nil
}

func (p *AccountsServiceDB) RegenerateToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	account_id := in.GetAccountId()
	fmt.Printf("Re-generating token for account_id = %d\n", account_id)

	account, err := getAndDecodeAccount(account_id, p.db)
	if err != nil {
		fmt.Printf("Failed to get token - %v", err)
		return nil, status.Error(codes.Internal, "Failed to Get token")
	}

	return &pb.TokenResponse{SessionToken: account.RegenerateToken()}, nil
}

func (p *AccountsServiceDB) GetAllAccounts(ctx context.Context, in *emptypb.Empty) (*pb.AccountListResponse, error) {
	fmt.Print("Listing all the accounts\n")

	account_ids, err := p.db.GetAllKeys()
	if err != nil {
		return nil, status.Error(codes.NotFound, "Failed to get accounts")
	}

	fmt.Printf("Accounts - %+v", account_ids)

	return &pb.AccountListResponse{AccountIds: account_ids}, nil
}

func getAndDecodeAccount(account_id int32, db *store.DbHandler) (*store.Account, error) {
	account_byte, err := db.Get(account_id)
	if err != nil {
		fmt.Printf("Failed to get account id = %d - error - %v", account_id, err)
		return nil, fmt.Errorf("failed to get account id - %d - error - %+v", account_id, err)
	}

	account, err := store.AccountDecode(account_byte)
	if err != nil {
		fmt.Printf("Failed to decode account - %d - error - %v", account_id, err)
		return nil, fmt.Errorf("failed to decode account id - %d - error - %+v", account_id, err)
	}
	return account, nil
}

type AccountsServiceMem struct {
	pb.UnimplementedAccountServer
	accounts_store *store.AccountsStore
	mu             sync.RWMutex
}

func NewAccountsServiceMem() *AccountsServiceMem {
	mem_store := store.NewAccountStore()
	return &AccountsServiceMem{
		accounts_store: mem_store,
	}
}

func (p *AccountsServiceMem) CreateAccount(ctx context.Context, in *pb.AccountInputRequest) (*pb.CreateAccountResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	account_id := in.GetAccountId()
	fmt.Printf("Creating account with account id - %d\n", account_id)
	res, err := p.accounts_store.CreateAccount(account_id)
	if err != nil {
		fmt.Printf("Failed to create account with account id - %d - error - %v", account_id, err)
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &pb.CreateAccountResponse{AccountId: account_id, SessionToken: res}, nil
}

func (p *AccountsServiceMem) GetToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	account_id := in.GetAccountId()
	fmt.Printf("Getting token for account_id = %d\n", account_id)
	res, err := p.accounts_store.GetToken(account_id)
	if err != nil {
		fmt.Printf("Failed to get token for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.TokenResponse{SessionToken: res}, nil
}

func (p *AccountsServiceMem) IsTokenExpired(ctx context.Context, in *pb.AccountInputRequest) (*pb.IsTokenExpiredResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	account_id := in.GetAccountId()
	fmt.Printf("Checking if token expired for account_id = %d\n", account_id)
	res, err := p.accounts_store.IsTokenExpired(account_id)
	if err != nil {
		fmt.Printf("Failed to check token expiry for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.Canceled, err.Error())
	}
	return &pb.IsTokenExpiredResponse{Expired: res}, nil
}

func (p *AccountsServiceMem) RegenerateToken(ctx context.Context, in *pb.AccountInputRequest) (*pb.TokenResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	account_id := in.GetAccountId()
	fmt.Printf("Re-generating token for account_id = %d\n", account_id)
	res, err := p.accounts_store.RegenerateToken(account_id)
	if err != nil {
		fmt.Printf("Failed to re-generate token for account id = %d - error - %v", account_id, err)
		return nil, status.Error(codes.Canceled, err.Error())
	}
	return &pb.TokenResponse{SessionToken: res}, nil
}

func (p *AccountsServiceMem) GetAllAccounts(ctx context.Context, in *emptypb.Empty) (*pb.AccountListResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	fmt.Print("Listing all the accounts\n")
	var account_ids []int32

	accounts, err := p.accounts_store.GetAccounts()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get accounts")
	}

	for key := range *accounts {
		account_ids = append(account_ids, key)
	}

	return &pb.AccountListResponse{AccountIds: account_ids}, nil
}
