package transaction

import "gorm.io/gorm"

type domain struct {
	ledgerRepo repoProvider
}

type repoProvider interface {
}

// New is to initialize ledger domain instance.
func New(db *gorm.DB) *domain {
	repo := newRepository(db)
	return &domain{
		ledgerRepo: repo,
	}
}
