package postgres

import (
	"app/internal/model"
	"context"
	"time"

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

	if err := r.db.WithContext(ctx).Table("deposits").Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *deposit) Get(ctx context.Context, depositId ksuid.KSUID) (*model.Deposit, error) {
	var out *model.Deposit
	err := r.db.WithContext(ctx).Table("deposits").Where("id = ?", depositId).Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *deposit) Update(ctx context.Context, in *model.Deposit) (*model.Deposit, error) {
	in.UpdatedAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Table("deposits").Save(in).Where("id = ?", in.Id).Error; err != nil {
		return nil, err
	}

	return in, nil
}
