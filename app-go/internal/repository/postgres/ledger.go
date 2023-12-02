package postgres

import (
	"app/internal/model"
	"context"

	"gorm.io/gorm"
)

type ledger struct {
	db *gorm.DB
}

func NewLedger(db *gorm.DB) *ledger {
	return &ledger{
		db: db,
	}
}

func (r *ledger) Create(ctx context.Context, entry *model.Ledger) (*model.Ledger, error) {
	return nil, nil
}
