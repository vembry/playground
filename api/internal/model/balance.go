package model

import (
	"errors"
	"time"

	"github.com/segmentio/ksuid"
)

var (
	ErrInsufficientBalance = errors.New("not enough balance")
	ErrBalanceLocked       = errors.New("balance is locked")
)

// Balance contain fields of balance. The fields refers to 'balances' table
type Balance struct {
	Id        ksuid.KSUID `json:"id" db:"id"`
	UserId    ksuid.KSUID `json:"user_id" db:"user_id"`
	Amount    float64     `json:"amount" db:"amount"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// WithdrawParam contain fields to execute balance withdrawal
type WithdrawParam struct {
	UserId      ksuid.KSUID `json:"user_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
}
