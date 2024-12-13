package postgres

import (
	"app/internal/model"
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type transfer struct {
	db *gorm.DB
}

func NewTransfer(db *gorm.DB) *transfer {
	return &transfer{
		db: db,
	}
}

func (r *transfer) Create(ctx context.Context, entry *model.Transfer) (*model.Transfer, error) {
	entry.Id = ksuid.New()

	if err := r.db.WithContext(ctx).Table("transfers").Create(entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *transfer) Get(ctx context.Context, transferId ksuid.KSUID) (*model.Transfer, error) {
	var out *model.Transfer
	err := r.db.WithContext(ctx).Table("transfers").Where("id = ?", transferId).Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *transfer) Update(ctx context.Context, in *model.Transfer) (*model.Transfer, error) {
	in.UpdatedAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Table("transfers").Save(in).Where("id = ?", in.Id).Error; err != nil {
		return nil, err
	}

	return in, nil
}
