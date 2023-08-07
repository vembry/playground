package repository

import "gorm.io/gorm"

type ledgerRepo struct {
	db *gorm.DB
}

// NewLedger is to initialize ledgers repository instance.
func NewLedger(db *gorm.DB) *ledgerRepo {
	return &ledgerRepo{
		db: db,
	}
}
