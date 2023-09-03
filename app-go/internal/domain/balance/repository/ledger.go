package repository

import (
	"app-go/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ledgerRepo struct {
	db *gorm.DB
}

// NewLedger is to initialize ledgers repository instance.
func NewLedger(db *gorm.DB) *ledgerRepo {
	return &ledgerRepo{
		db: db,
	}
}

// Create is to create new ledger entry
func (lr *ledgerRepo) Create(ctx context.Context, in *model.Ledger) error {
	res := lr.db.WithContext(ctx).Table("ledgers").Save(in)
	if res.Error != nil {
		return fmt.Errorf("found error on inserting ledger to db. err=%w", res.Error)
	}
	return nil
}
