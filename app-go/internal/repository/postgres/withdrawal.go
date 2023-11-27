package postgres

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type withdrawal struct {
	db *gorm.DB
}

func NewWithdrawal(db *gorm.DB) *withdrawal {
	return &withdrawal{
		db: db,
	}
}

func (r *withdrawal) Create(ctx context.Context, entry *model.Withdrawal) (*model.Withdrawal, error) {
	entry.Id = ksuid.New()

	if err := r.db.Table("withdrawals").Create(entry).WithContext(ctx).Error; err != nil {
		return nil, err
	}
	return entry, nil
}
