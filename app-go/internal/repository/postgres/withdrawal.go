package postgres

import (
	"app/internal/model"
	"context"
	"time"

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

func (r *withdrawal) Get(ctx context.Context, withdrawalId ksuid.KSUID) (*model.Withdrawal, error) {
	var out *model.Withdrawal
	err := r.db.Table("withdrawals").Where("id = ?", withdrawalId).Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *withdrawal) Update(ctx context.Context, in *model.Withdrawal) (*model.Withdrawal, error) {
	in.UpdatedAt = time.Now().UTC()
	if err := r.db.Table("withdrawals").Save(in).Where("id = ?", in.Id).Error; err != nil {
		return nil, err
	}

	return in, nil
}
