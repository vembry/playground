package transaction

import (
	"api/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// repository is transaction's repository instance
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

// Get is to get a single data from `transactions` table by transactionId
func (r *repository) Get(ctx context.Context, transactionId ksuid.KSUID) (*model.Transaction, error) {
	var out model.Transaction
	res := r.db.WithContext(ctx).Table("transactions").First(&out, "id = ?", transactionId.String())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("found error on getting transaction from db. transactionId=%s. err=%w", transactionId, res.Error)
	}

	return &out, nil
}

// Update is to update existing transaction data
func (r *repository) Update(ctx context.Context, in *model.Transaction) error {
	res := r.db.WithContext(ctx).Table("transactions").Save(in)
	if res.Error != nil {
		return fmt.Errorf("found error on updating transaction to db. err=%w", res.Error)
	}
	return nil
}
