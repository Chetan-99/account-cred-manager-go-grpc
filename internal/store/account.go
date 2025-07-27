package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"time"
)

const (
	SESSION_TOKEN_EXPIRY_SECONDS = 60 * 9
)

type Account struct {
	AccountId int32

	SessionToken          string
	CreatedSessionTokenTS time.Time
}

func NewAccount(account_id int32) *Account {
	token, created_ts := createToken()

	return &Account{
		AccountId:             account_id,
		SessionToken:          token,
		CreatedSessionTokenTS: created_ts,
	}
}

func AccountDecode(data []byte) (*Account, error) {
	var ta Account
	account_buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(account_buf)
	if err := decoder.Decode(&ta); err != nil {
		return &Account{}, fmt.Errorf("failed to decode account - error - %+v", err)
	}
	return &ta, nil
}

func createToken() (string, time.Time) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	token := fmt.Sprintf("%d", r.Intn(900000)+100000)
	return token, time.Now()
}

func (a *Account) Print() {
	println("account_id = %d", a.AccountId)
	println("session_token = %s", a.SessionToken)
	println("created ts = %s", a.CreatedSessionTokenTS.String())
}

func (a *Account) IsTokenExpired() bool {
	elapsed := time.Since(a.CreatedSessionTokenTS).Seconds()
	return elapsed > SESSION_TOKEN_EXPIRY_SECONDS
}

func (a *Account) RegenerateToken() string {
	token, created_ts := createToken()
	a.SessionToken = token
	a.CreatedSessionTokenTS = created_ts
	return token
}

func (a *Account) GetToken() string {
	if a.IsTokenExpired() {
		return a.RegenerateToken()
	} else {
		return a.SessionToken
	}
}

func (a *Account) Encode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(*a)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
