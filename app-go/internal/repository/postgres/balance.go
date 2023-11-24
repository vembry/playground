package postgres

import (
	"app/internal/model"
	"context"

	"gorm.io/gorm"
)

type balance struct {
	db *gorm.DB
}

func NewBalance(db *gorm.DB) *balance {
	return &balance{
		db: db,
	}
}
func (r *balance) Create(ctx context.Context, entry *model.Balance) (*model.Balance, error) {
	return nil, nil
}

func (r *balance) Update(ctx context.Context, in *model.Balance) (*model.Balance, error) {
	return nil, nil
}
