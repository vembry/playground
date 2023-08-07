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
	transactionRepo           repoProvider
	ledger                    ledgerProvider
	pendingTransactionHandler pendingTransactionHandlerProvider
}

// repoProvider is the spec of transaction's repository
type repoProvider interface {
	Create(ctx context.Context, in *model.Transaction) error
}

// ledgerProvider is the spec of ledger's instance
type ledgerProvider interface {
	CreateEntry(ctx context.Context, in *model.CreateLedgerEntry) error
}

// pendingTransactionHandlerProvider is the spec of transcation-pending worker's handler
type pendingTransactionHandlerProvider interface {
	Enqueue(ctx context.Context, transactionId ksuid.KSUID) error
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

// WithPendingTransactionHandler is to inject ledger into pending-transaction handler
func (d *domain) WithPendingTransactionHandler(pendingTransactionHandler pendingTransactionHandlerProvider) {
	d.pendingTransactionHandler = pendingTransactionHandler
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

	// create transaction entry
	err = d.transactionRepo.Create(ctx, &newTransaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction. err=%w", err)
	}

	// enqueue transaction to the worker
	err = d.pendingTransactionHandler.Enqueue(ctx, newTransaction.Id)
	if err != nil {
		return fmt.Errorf("failed to enqueue pending-transaction. transaction-id=%s. err=%w", newTransaction.Id, err)
	}

	return nil
}

// ProcessPending is to proceed pending transactions
func (d *domain) ProcessPending(ctx context.Context, transactionId ksuid.KSUID) error {
	return nil
}
