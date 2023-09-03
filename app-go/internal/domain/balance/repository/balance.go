package repository

import (
	"app-go/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type balanceRepo struct {
	db *model.DB
}

// newRepository is to initialize balances repository instance.
func NewBalance(db *model.DB) *balanceRepo {
	return &balanceRepo{
		db: db,
	}
}

// Get is to get balance by userId
func (br *balanceRepo) Get(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error) {
	var out model.Balance
	res := br.db.Slave.WithContext(ctx).Table("balances").First(&out, "user_id = ?", userId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("found error on retrieving balances from db. err=%w", res.Error)
	}
	return &out, nil
}

// Get is to get balance by userId from master db
func (br *balanceRepo) GetFromMaster(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error) {
	var out model.Balance
	res := br.db.Master.WithContext(ctx).Table("balances").First(&out, "user_id = ?", userId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("found error on retrieving balances from db. err=%w", res.Error)
	}
	return &out, nil
}

// Update is to update existing balance data
func (br *balanceRepo) Update(ctx context.Context, balance *model.Balance) error {
	res := br.db.Master.WithContext(ctx).Table("balances").Save(balance)
	if res.Error != nil {
		return fmt.Errorf("found error on updating balances to db. err=%w", res.Error)
	}
	return nil
}
