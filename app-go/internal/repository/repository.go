package repository

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

type IBalance interface {
	Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error)
	Create(ctx context.Context, entry *model.Balance) (*model.Balance, error)
	Update(ctx context.Context, in *model.Balance) (*model.Balance, error)
}

type ILedger interface {
	Create(ctx context.Context, entry *model.Ledger) (*model.Ledger, error)
}

type IDeposit interface {
	Create(ctx context.Context, entry *model.Deposit) (*model.Deposit, error)
}
type IWithdrawal interface {
	Create(ctx context.Context, entry *model.Withdrawal) (*model.Withdrawal, error)
}
type ITransfer interface {
	Create(ctx context.Context, entry *model.Transfer) (*model.Transfer, error)
}
