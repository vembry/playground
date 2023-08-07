package repository

import "gorm.io/gorm"

type balanceRepo struct {
	db *gorm.DB
}

// newRepository is to initialize balances repository instance.
func NewBalance(db *gorm.DB) *balanceRepo {
	return &balanceRepo{
		db: db,
	}
}
