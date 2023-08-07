package transaction

import (
	"api/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// domain is transaction's domain instance
type domain struct {
	transactionRepo repoProvider
	ledger          ledgerProvider
}

// repoProvider is the spec of transaction's repository
type repoProvider interface {
	Create(ctx context.Context, in *model.Transaction) error
}

// ledgerProvider is the spec of ledger's instance
type ledgerProvider interface {
	CreateEntry(ctx context.Context, in *model.CreateLedgerEntry) error
}

// New is to initialize transaction domain instance.
func New(db *gorm.DB) *domain {
	repo := newRepository(db)
	return &domain{
		transactionRepo: repo,
	}
}

// WithLedger is to inject ledger into transaction domain
func (d *domain) WithLedger(ledger ledgerProvider) {
	d.ledger = ledger
}

// Create is to create a single transaction entry
func (d *domain) Create(ctx context.Context, in *model.CreateTransaction) error {
	userId, err := ksuid.Parse(in.UserId)
	if err != nil {
		return fmt.Errorf("failed parsing user-id. err=%w", err)
	}

	newTransaction := model.Transaction{
		Id:          ksuid.New(),
		UserId:      userId,
		Status:      model.TransactionStatusPending,
		Description: in.Description,
		Remarks:     "",
		Amount:      in.Amount,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	err = d.transactionRepo.Create(ctx, &newTransaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction. err=%w", err)
	}

	return nil
}
