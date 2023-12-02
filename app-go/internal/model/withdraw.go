package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Withdrawal struct {
	Id        ksuid.KSUID `json:"id" gorm:"column:id"`
	BalanceId ksuid.KSUID `json:"balance_id" gorm:"column:balance_id"`
	Status    Status      `json:"status" gorm:"column:status"`
	Amount    float64     `json:"amount" gorm:"column:amount"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"column:updated_at"`
}
