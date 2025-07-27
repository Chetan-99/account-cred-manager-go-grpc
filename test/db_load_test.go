package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chetan-99/account-cred-manager-go-grpc/internal/config"
	"github.com/chetan-99/account-cred-manager-go-grpc/internal/store"
)

var testDb *store.DbHandler

func TestMain(m *testing.M) {

	// Pre-test steps

	cfg := config.AppConfig{
		DB_PATH: "./tmp.db",
	}

	os.RemoveAll("./tmp.db")

	testDb = store.NewBadgerDB(&cfg)
	defer testDb.Close()

	// Tests

	code := m.Run()

	// Post-test steps
	err := os.RemoveAll("./tmp.db")
	if err != nil {
		fmt.Printf("failed to removed temporary test db")
	}

	os.Exit(code)
}

func TestDBLoadFunctionTesting(t *testing.T) {
	if testDb == nil {
		t.Fatal("testDb is nil")
	}

	// Create test data
	testAccountIds := [5]int32{12, 34, 54, 21, 85}
	testAccounts := []*store.Account{}

	for _, account_id := range testAccountIds {
		testAccounts = append(testAccounts, store.NewAccount(account_id))
	}

	// Insert test accounts into DB
	for _, acc := range testAccounts {
		data, err := acc.Encode()
		if err != nil {
			t.Fatalf("failed to encode account: %v", err)
		}
		err = testDb.Add_KV(acc.AccountId, data)
		if err != nil {
			t.Fatalf("failed to insert account: %v", err)
		}
	}

	// Load accounts from DB and verify
	var loadedAccounts []*store.Account
	for _, acc := range testAccounts {
		data, err := testDb.Get(acc.AccountId)
		if err != nil {
			t.Fatalf("failed to get account: %v", err)
		}
		loadedAcc, err := store.AccountDecode(data)
		if err != nil {
			t.Fatalf("failed to decode account: %v", err)
		}
		loadedAccounts = append(loadedAccounts, loadedAcc)
	}
	if len(loadedAccounts) != len(testAccounts) {
		t.Fatalf("expected %d accounts, got %d", len(testAccounts), len(loadedAccounts))
	}

	// Check account data matches
	for i, acc := range loadedAccounts {
		if acc.AccountId != testAccounts[i].AccountId || acc.SessionToken != testAccounts[i].SessionToken {
			t.Errorf("account mismatch: got %+v, want %+v", acc, testAccounts[i])
		}
	}

	t.Logf("Completed - DB Load Function, loaded %d accounts", len(loadedAccounts))
}
