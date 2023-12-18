package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Ledger struct {
	Id            ksuid.KSUID `json:"id" gorm:"column:id"`
	BalanceId     ksuid.KSUID `json:"balance_id" gorm:"column:balance_id"`
	Type          LedgerType  `json:"type" gorm:"column:type"`
	Amount        float64     `json:"amount" gorm:"column:amount"`
	BalanceAfter  float64     `json:"balance_after" gorm:"column:balance_after"`
	BalanceBefore float64     `json:"balance_before" gorm:"column:balance_before"`
	CreatedAt     time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"column:updated_at"`
}

type LedgerType string

const (
	LedgerTypeIn  LedgerType = "in"
	LedgerTypeOut LedgerType = "out"
)
