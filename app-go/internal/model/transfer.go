package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Transfer struct {
	Id            ksuid.KSUID `json:"id" gorm:"column:id"`
	BalanceIdFrom ksuid.KSUID `json:"balance_id_from" gorm:"column:balance_id_from"`
	BalanceIdTo   ksuid.KSUID `json:"balance_id_to" gorm:"column:balance_id_to"`
	Status        Status      `json:"status" gorm:"column:status"`
	Amount        float64     `json:"amount" gorm:"column:amount"`
	CreatedAt     time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"column:updated_at"`
}
