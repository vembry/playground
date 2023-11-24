package postgres

import (
	"app/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/segmentio/ksuid"
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
	// assign basic values
	entry.Id = ksuid.New()
	entry.CreatedAt = time.Now().UTC()
	entry.UpdatedAt = entry.CreatedAt

	// insert to table
	if err := r.db.Table("balances").Create(entry).Error; err != nil {
		if err.(*pgconn.PgError).Code == "23505" {
			return nil, fmt.Errorf("balance exists")
		}
		return nil, err
	}
	return entry, nil
}

func (r *balance) Update(ctx context.Context, in *model.Balance) (*model.Balance, error) {
	in.UpdatedAt = time.Now().UTC()

	if err := r.db.Table("balances").Save(in).Where("id = ?", in.Id).Error; err != nil {
		return nil, err
	}

	return in, nil
}
