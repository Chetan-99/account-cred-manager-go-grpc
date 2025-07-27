package store

import (
	"errors"
	"sync"
)

type AccountsStore struct {
	accounts map[int32]Account
	mu       sync.RWMutex
}

func NewAccountStore() *AccountsStore {
	return &AccountsStore{
		accounts: make(map[int32]Account),
	}
}

func (a *AccountsStore) GetAccounts() (*map[int32]Account, error) {
	return &a.accounts, nil
}

func (a *AccountsStore) CreateAccount(account_id int32) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.accounts[account_id]; ok {
		return "", errors.New("ALREADY_EXIST")
	} else {
		a.accounts[account_id] = *NewAccount(account_id)
		return a.accounts[account_id].SessionToken, nil
	}
}

func (a *AccountsStore) IsTokenExpired(account_id int32) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if val, ok := a.accounts[account_id]; ok {
		return val.IsTokenExpired(), nil
	} else {
		return false, errors.New("ACCOUNT_DOES_NOT_EXIST")
	}
}

func (a *AccountsStore) RegenerateToken(account_id int32) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if val, ok := a.accounts[account_id]; ok {
		return val.RegenerateToken(), nil
	} else {
		return "", errors.New("ACCOUNT_DOES_NOT_EXIST")
	}
}

func (a *AccountsStore) GetToken(account_id int32) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if val, ok := a.accounts[account_id]; ok {
		return val.GetToken(), nil
	} else {
		return "", errors.New("ACCOUNT_DOES_NOT_EXIST")
	}
}
