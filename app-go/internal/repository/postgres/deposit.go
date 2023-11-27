package postgres

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type deposit struct {
	db *gorm.DB
}

func NewDeposit(db *gorm.DB) *deposit {
	return &deposit{
		db: db,
	}
}

func (r *deposit) Create(ctx context.Context, entry *model.Deposit) (*model.Deposit, error) {
	entry.Id = ksuid.New()

	if err := r.db.Table("deposits").Create(entry).WithContext(ctx).Error; err != nil {
		return nil, err
	}
	return entry, nil
}
