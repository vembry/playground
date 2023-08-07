package transaction

import (
	"api/internal/domain/ledger/repository"
	"api/internal/model"
	"context"

	"gorm.io/gorm"
)

// domain is ledger's domain instance
type domain struct {
	ledgerRepo  ledgerRepoProvider
	balanceRepo balanceRepoProvider
}

// ledgerRepoProvider is the spec of ledger's repository
type ledgerRepoProvider interface {
}

// balanceRepoProvider is the spec of balance's repository
type balanceRepoProvider interface {
}

// New is to initialize ledger domain instance.
func New(db *gorm.DB) *domain {
	ledgerRepo := repository.NewLedger(db)
	balanceRepo := repository.NewBalance(db)

	return &domain{
		ledgerRepo:  ledgerRepo,
		balanceRepo: balanceRepo,
	}
}

// CreateEntry is to create ledger entry
func (d *domain) CreateEntry(ctx context.Context, in *model.CreateLedgerEntry) error {
	return nil
}

// GetBalance is to get user's balance
func (d *domain) GetBalance(ctx context.Context, userId string) error {
	return nil
}
