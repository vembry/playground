package transaction

import (
	"api/internal/domain/balance/repository"
	"api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// domain is balance's domain instance
type domain struct {
	ledgerRepo  ledgerRepoProvider
	balanceRepo balanceRepoProvider
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

// New is to initialize balance domain instance.
func New(db *gorm.DB) *domain {
	ledgerRepo := repository.NewLedger(db)
	balanceRepo := repository.NewBalance(db)

	return &domain{
		ledgerRepo:  ledgerRepo,
		balanceRepo: balanceRepo,
	}
}

// Withdraw is to withdraw credit from balance
func (d *domain) Withdraw(ctx context.Context, in *model.WithdrawParam) error {
	// get balance
	balance, err := d.GetBalance(ctx, in.UserId)
	if err != nil {
		return fmt.Errorf("found error on getting balance by userId. userId=%s. err=%w", in.UserId, err)
	}

	// validate amount withdrawn against amount available on current balance
	if balance.Amount < in.Amount {
		return fmt.Errorf("not enough balance. balanceId=%s", balance.Id)
	}

	// transform balance
	newBalance := *balance

	// update balance values
	newBalance.Amount -= in.Amount
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

// GetBalance is to get user's balance
func (d *domain) GetBalance(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error) {
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
