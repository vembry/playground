package repository

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

type IBalance interface {
	Create(ctx context.Context, entry *model.Balance) (*model.Balance, error)
	Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error)
	Update(ctx context.Context, in *model.Balance) (*model.Balance, error)
}

type ILedger interface {
	Create(ctx context.Context, entry *model.Ledger) (*model.Ledger, error)
}

type IDeposit interface {
	Create(ctx context.Context, entry *model.Deposit) (*model.Deposit, error)
	Get(ctx context.Context, depositId ksuid.KSUID) (*model.Deposit, error)
	Update(ctx context.Context, in *model.Deposit) (*model.Deposit, error)
}
type IWithdrawal interface {
	Create(ctx context.Context, entry *model.Withdrawal) (*model.Withdrawal, error)
	Get(ctx context.Context, withdrawalId ksuid.KSUID) (*model.Withdrawal, error)
	Update(ctx context.Context, in *model.Withdrawal) (*model.Withdrawal, error)
}
type ITransfer interface {
	Create(ctx context.Context, entry *model.Transfer) (*model.Transfer, error)
	Get(ctx context.Context, transferId ksuid.KSUID) (*model.Transfer, error)
	Update(ctx context.Context, in *model.Transfer) (*model.Transfer, error)
}
