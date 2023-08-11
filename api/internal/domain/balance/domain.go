package transaction

import (
	"api/internal/domain/balance/repository"
	"api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// domain is balance's domain instance
type domain struct {
	ledgerRepo  ledgerRepoProvider
	balanceRepo balanceRepoProvider
	mutex       mutexProvider
}

// ledgerRepoProvider is the spec of ledger's repository
type ledgerRepoProvider interface {
	Create(ctx context.Context, in *model.Ledger) error
}

// balanceRepoProvider is the spec of balance's repository
type balanceRepoProvider interface {
	Get(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error)
	Update(ctx context.Context, balance *model.Balance) error
}

// mutexProvider is the spec of mutex's instance
type mutexProvider interface {
	Acquire(ctx context.Context, name string) (*model.Mutex, error)
	Delete(ctx context.Context, mutex *model.Mutex) error
}

// New is to initialize balance domain instance.
func New(db *gorm.DB, mutex mutexProvider) *domain {
	ledgerRepo := repository.NewLedger(db)
	balanceRepo := repository.NewBalance(db)

	return &domain{
		ledgerRepo:  ledgerRepo,
		balanceRepo: balanceRepo,
		mutex:       mutex,
	}
}

// Withdraw is to withdraw credit from balance
func (d *domain) Withdraw(ctx context.Context, in *model.WithdrawBalanceParam) error {
	// get balance lock
	balance, unlockBalance, errLocked, err := d.GetLock(ctx, in.UserId)
	if errLocked != nil {
		return errLocked
	}
	if err != nil {
		return fmt.Errorf("found error on getting balance lock by userId. userId=%s. err=%w", in.UserId.String(), err)
	}

	defer unlockBalance(ctx)

	// validate amount withdrawn against amount available on current balance
	if balance.Amount < in.Amount {
		return model.ErrInsufficientBalance
	}

	// transform balance
	newBalance := *balance

	a := big.NewFloat(balance.Amount)
	b := big.NewFloat(in.Amount)

	c := new(big.Float).Sub(a, b)

	// update balance values
	newBalance.Amount, _ = c.Float64()
	newBalance.UpdatedAt = time.Now().UTC()

	// save updated balance
	err = d.balanceRepo.Update(ctx, &newBalance)
	if err != nil {
		return fmt.Errorf("found error on updating balance. balanceId=%s. err=%w", balance.Id, err)
	}

	// create ledger entry
	newLedgerEntry := model.Ledger{
		Id:            ksuid.New(),
		UserId:        balance.UserId,
		Type:          model.LedgerTypeOut,
		Description:   in.Description,
		Amount:        in.Amount,
		BalanceAfter:  newBalance.Amount,
		BalanceBefore: balance.Amount,
		CreatedAt:     time.Now().UTC(),
	}

	// save new ledger entry
	err = d.ledgerRepo.Create(ctx, &newLedgerEntry)
	if err != nil {
		raw, _ := json.Marshal(newLedgerEntry)
		log.Printf("found error on creating ledger. ledgerEntry=%s. err=%v", string(raw), err)
	}

	return nil
}

// Get is to get user's balance
func (d *domain) Get(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error) {
	// get balance
	balance, err := d.balanceRepo.Get(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("found error on getting balance by userId. userId=%s. err=%w", userId, err)
	}

	// validate balance existence
	if balance == nil {
		return nil, fmt.Errorf("balance not found. userId=%s", userId)
	}
	return balance, nil
}

// GetLock is to get and lock user's balance. Return balance, balance-unlocker, locked-error, and generic-error
func (d *domain) GetLock(ctx context.Context, userId ksuid.KSUID) (*model.Balance, func(context.Context), error, error) {
	// get balance
	balance, err := d.Get(ctx, userId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("found error on getting balance by userId. userId=%s. err=%w", userId, err)
	}

	//  lock balance
	mutexRes, err := d.mutex.Acquire(ctx, fmt.Sprintf("balanceLocker.%s", balance.Id))
	_ = mutexRes
	if err != nil {
		return nil, nil, model.ErrBalanceLocked, nil
	}

	return balance, func(_ctx context.Context) {
		// balance unlock
		if err := d.mutex.Delete(_ctx, mutexRes); err != nil {
			log.Printf("found error on deleting mutex. mutexName=%s. err=%v", mutexRes.Name, err)
		}
	}, nil, nil
}

// Add is to add credit to active balance
func (d *domain) Add(ctx context.Context, in *model.AddBalanceParam) error {
	// get balance lock
	balance, unlockBalance, errLocked, err := d.GetLock(ctx, in.UserId)
	if errLocked != nil {
		return errLocked
	}
	if err != nil {
		return fmt.Errorf("found error on getting balance lock by userId. userId=%s. err=%w", in.UserId, err)
	}

	defer unlockBalance(ctx)

	// transform balance
	newBalance := *balance

	a := big.NewFloat(balance.Amount)
	b := big.NewFloat(in.Amount)

	c := new(big.Float).Add(a, b)

	// update balance values
	newBalance.Amount, _ = c.Float64()
	newBalance.UpdatedAt = time.Now().UTC()

	// save updated balance
	err = d.balanceRepo.Update(ctx, &newBalance)
	if err != nil {
		return fmt.Errorf("found error on updating balance. balanceId=%s. err=%w", balance.Id, err)
	}

	// create ledger entry
	newLedgerEntry := model.Ledger{
		Id:            ksuid.New(),
		UserId:        balance.UserId,
		Type:          model.LedgerTypeIn,
		Description:   "topup",
		Amount:        in.Amount,
		BalanceAfter:  newBalance.Amount,
		BalanceBefore: balance.Amount,
		CreatedAt:     time.Now().UTC(),
	}

	// save new ledger entry
	err = d.ledgerRepo.Create(ctx, &newLedgerEntry)
	if err != nil {
		raw, _ := json.Marshal(newLedgerEntry)
		log.Printf("found error on creating ledger. ledgerEntry=%s. err=%v", string(raw), err)
	}

	return nil
}
