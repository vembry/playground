package transaction

import (
	"api/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// domain is transaction's domain instance
type domain struct {
	transactionRepo           repoProvider
	balance                   balanceProvider
	pendingTransactionHandler pendingTransactionHandlerProvider
}

// repoProvider is the spec of transaction's repository
type repoProvider interface {
	Create(ctx context.Context, in *model.Transaction) error
	Get(ctx context.Context, transactionId ksuid.KSUID) (*model.Transaction, error)
}

// balanceProvider is the spec of balance's instance
type balanceProvider interface {
	Withdraw(ctx context.Context, in *model.WithdrawParam) error
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

// WithBalance is to inject balance into transaction domain
func (d *domain) WithBalance(balance balanceProvider) {
	d.balance = balance
}

// WithPendingTransactionHandler is to inject pending-transaction handler dependencies
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
		return fmt.Errorf("failed to enqueue pending-transaction. transactionId=%s. err=%w", newTransaction.Id, err)
	}

	return nil
}

// ProcessPending is to proceed pending transactions
func (d *domain) ProcessPending(ctx context.Context, transactionId ksuid.KSUID) error {
	log.Printf("processing pending transaction. transaction-id=%s", transactionId)

	// get transaction by transactionId
	transaction, err := d.transactionRepo.Get(ctx, transactionId)
	if err != nil {
		return fmt.Errorf("found error on getting transaction by id. transactionId=%s. err=%w", transactionId, err)
	}

	// withdraw money
	err = d.balance.Withdraw(ctx, &model.WithdrawParam{
		UserId:      transaction.UserId,
		Amount:      transaction.Amount,
		Description: transaction.Description,
	})
	if err != nil {
		return fmt.Errorf("found error on withdrawing money from balance. transactionId=%s. err=%w", transactionId, err)
	}

	return nil
}
