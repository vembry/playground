package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

// CreateTransaction contain fields to create a transaction entry
type CreateTransaction struct {
	UserId      ksuid.KSUID `json:"user_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
}

// Transaction contain fields of transactions. The fields refers to 'transactions' table
type Transaction struct {
	Id          ksuid.KSUID `json:"id" db:"id"`
	UserId      ksuid.KSUID `json:"user_id" db:"user_id"`
	Status      string      `json:"status" db:"status"`
	Description string      `json:"description" db:"description"`
	Remarks     string      `json:"remarks" db:"remarks"`
	Amount      float64     `json:"amount" db:"amount"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// containing transaction's status enum
const (
	TransactionStatusPending string = "pending"
	TransactionStatusSuccess string = "success"
	TransactionStatusFailed  string = "failed"
)
