package transaction

import (
	"api/internal/model"
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// newRepository is to initialize transactions repository instance.
func newRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

// Create is to create an entry to the `transactions` table
func (r *repository) Create(ctx context.Context, in *model.Transaction) error {
	res := r.db.Create(in)
	return res.Error
}
