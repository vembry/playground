package postgres

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
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
	entry.Id = ksuid.New()
	if err := r.db.WithContext(ctx).Table("ledgers").Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}
