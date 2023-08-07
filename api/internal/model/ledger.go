package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

// CreateLedgerEntry contain fields to create a ledger entry
type CreateLedgerEntry struct {
	UserId      ksuid.KSUID `json:"user_id" db:"user_id"`
	Type        string      `json:"type" db:"type"`
	Description string      `json:"description" db:"description"`
	Amount      float64     `json:"amount" db:"amount"`
}

// Ledger contain fields of ledger. The fields refers to 'ledgers' table
type Ledger struct {
	Id            ksuid.KSUID `json:"id" db:"id"`
	UserId        ksuid.KSUID `json:"user_id" db:"user_id"`
	Type          string      `json:"type" db:"type"`
	Description   string      `json:"description" db:"description"`
	Amount        float64     `json:"amount" db:"amount"`
	BalanceAfter  float64     `json:"balance_after" db:"balance_after"`
	BalanceBefore float64     `json:"balance_before" db:"balance_before"`
	CreatedAt     time.Time   `json:"create_at" db:"created_at"`
}

const (
	LedgerTypeIn  string = "in"
	LedgerTypeOut string = "out"
)
