package domain

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
)

type IBalance interface {
	Open(ctx context.Context) (*model.Balance, error)
	Deposit(ctx context.Context, in *model.DepositParam) error
	Withdraw(ctx context.Context, in *model.WithdrawParam) error
	Transfer(ctx context.Context, in *model.TransferParam) error
}

type balance struct {
	balanceRepository repository.IBalance
	ledgerRepository  repository.ILedger
}

func NewBalance(
	balanceRepository repository.IBalance,
	ledgerRepository repository.ILedger,
) *balance {
	return &balance{
		balanceRepository: balanceRepository,
		ledgerRepository:  ledgerRepository,
	}
}

func (d *balance) Open(ctx context.Context) (*model.Balance, error) {
	return nil, nil
}

func (d *balance) Deposit(ctx context.Context, in *model.DepositParam) error {
	return nil
}

func (d *balance) Withdraw(ctx context.Context, in *model.WithdrawParam) error {
	return nil
}

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) error {
	return nil
}
