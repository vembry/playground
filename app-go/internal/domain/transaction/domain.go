package transaction

import (
	"app/internal/model"
	"context"
	"encoding/json"
	"errors"
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
	Update(ctx context.Context, in *model.Transaction) error
}

// balanceProvider is the spec of balance's instance
type balanceProvider interface {
	Withdraw(ctx context.Context, in *model.WithdrawBalanceParam) error
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
	newTransaction := model.Transaction{
		Id:          ksuid.New(),
		UserId:      in.UserId,
		Status:      model.TransactionStatusPending,
		Description: in.Description,
		Remarks:     "",
		Amount:      in.Amount,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// create transaction entry
	err := d.transactionRepo.Create(ctx, &newTransaction)
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

	finalTransactionStatus := model.TransactionStatusSuccess

	// withdraw money
	err = d.balance.Withdraw(ctx, &model.WithdrawBalanceParam{
		UserId:      transaction.UserId,
		Amount:      transaction.Amount,
		Description: transaction.Description,
	})
	if err != nil {
		if errors.Is(err, model.ErrInsufficientBalance) {
			// for insufficient balance, we want to stop the queue right away
			log.Printf("not enough balance to do transaction. transactionId=%s", transaction.Id)
			// return nil
			finalTransactionStatus = model.TransactionStatusFailed
		} else if errors.Is(err, model.ErrBalanceLocked) {
			// for locked balance we want to requeue the transaction
			// instead of waiting for the balance get unlocked
			// and then do early return
			log.Printf("balance is currently locked, transaction requeued. transactionId=%s", transaction.Id)
			d.pendingTransactionHandler.Enqueue(ctx, transaction.Id)
			return nil
		} else {
			// else, just return the error wrapped
			return fmt.Errorf("found error on withdrawing money from balance. transactionId=%s. err=%w", transactionId, err)
		}
	}

	// update transaction
	newTransaction := *transaction
	newTransaction.Status = finalTransactionStatus

	// save updated transaction
	err = d.transactionRepo.Update(ctx, &newTransaction)
	if err != nil {
		// if there is error, dont break the process
		// but leave logs. or else
		raw, _ := json.Marshal(newTransaction)
		log.Printf("found error on updating transaction. transactionId=%s. transaction=%s. err=%v", transaction.Id, string(raw), err)
	}

	return nil
}
