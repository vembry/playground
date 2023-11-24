package repository

import (
	"app/internal/model"
	"context"
)

type IBalance interface {
	Create(ctx context.Context, entry *model.Balance) (*model.Balance, error)
	Update(ctx context.Context, in *model.Balance) (*model.Balance, error)
}

type ILedger interface {
	Create(ctx context.Context, entry *model.Ledger) (*model.Ledger, error)
}
