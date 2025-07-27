package store

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	SESSION_TOKEN_EXPIRY_SECONDS = 60 * 9
)

type account struct {
	account_id int32

	session_token            string
	created_session_token_ts time.Time
}

func newAccount(account_id int32) *account {
	token, created_ts := createToken()

	return &account{
		account_id:               account_id,
		session_token:            token,
		created_session_token_ts: created_ts,
	}
}

func createToken() (string, time.Time) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	token := fmt.Sprintf("%d", r.Intn(900000)+100000)
	return token, time.Now()
}

func (a *account) IsTokenExpired() bool {
	elapsed := time.Since(a.created_session_token_ts).Seconds()
	return elapsed > SESSION_TOKEN_EXPIRY_SECONDS
}

func (a *account) RegenerateToken() string {
	token, created_ts := createToken()
	a.session_token = token
	a.created_session_token_ts = created_ts
	return token
}

func (a *account) GetToken() string {
	if a.IsTokenExpired() {
		return a.RegenerateToken()
	} else {
		return a.session_token
	}
}
