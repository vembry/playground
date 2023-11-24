package model

import "github.com/segmentio/ksuid"

type Balance struct {
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
