package balance

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessDeposit(ctx context.Context, depositId ksuid.KSUID) error {
	// get deposit data
	deposit, err := d.depositRepo.Get(ctx, depositId)
	if err != nil {
		return err
	}

	// validate state
	if deposit == nil || deposit.Status != model.StatusPending {
		return nil
	}

	// get balance lock
	balance, unlocker, err := d.GetLock(ctx, deposit.BalanceId)
	if unlocker != nil {
		defer unlocker(ctx)
	}
	if err != nil {
		// produce task for worker
		d.depositProducer.Produce(ctx, depositId)
		return nil
	}

	// construct ledger entry
	ledger := &model.Ledger{
		BalanceId:     balance.Id,
		Type:          model.LedgerTypeIn,
		Amount:        deposit.Amount,
		BalanceBefore: balance.Amount,
	}

	// add balance
	balance.Amount += deposit.Amount

	// assign new balance
	ledger.BalanceAfter = balance.Amount

	// update balance
	_, err = d.balanceRepo.Update(ctx, balance)
	if err != nil {
		return err
	}

	deposit.Status = model.StatusCompleted

	// update withdrawal
	_, err = d.depositRepo.Update(ctx, deposit)
	if err != nil {
		return err
	}

	// create ledger entry
	ledger, err = d.ledgerRepo.Create(ctx, ledger)
	if err != nil {
		return err
	}

	return nil
}
