package auth

import "errors"

const ProviderFirebase = "firebase"

var (
	ErrNotFound  = errors.New("account not synced")
	ErrSuspended = errors.New("account suspended")
)

type Identity struct {
	Provider string
	Subject  string
}

type Account struct {
	ID      string `json:"account_id"`
	Status  string `json:"status"`
	Created bool   `json:"-"`
}
