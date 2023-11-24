package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type Balance struct {
	Id        ksuid.KSUID `json:"id" gorm:"column:id"`
	Amount    float64     `json:"amount" gorm:"column:amount"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"column:updated_at"`
}

type DepositParam struct {
	BalanceId ksuid.KSUID `json:"balance_id"`
	Amount    float64     `json:"amount"`
}

type WithdrawParam struct {
	BalanceId ksuid.KSUID `json:"balance_id"`
	Amount    float64     `json:"amount"`
}

type TransferParam struct {
	BalanceIdFrom ksuid.KSUID `json:"balance_id_from"`
	BalanceIdTo   ksuid.KSUID `json:"balance_id_to"`
	Amount        float64     `json:"amount"`
}
