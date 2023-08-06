package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

// newRepository is to initialize ledgers repository instance.
func newRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}
