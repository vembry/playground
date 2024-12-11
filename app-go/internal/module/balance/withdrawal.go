package balance

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get withdraw
	withdrawal, err := d.withdrawalRepo.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	// validate state
	if withdrawal == nil || withdrawal.Status != model.StatusPending {
		return nil
	}

	// get balance lock
	balance, unlocker, err := d.GetLock(ctx, withdrawal.BalanceId)
	if unlocker != nil {
		defer unlocker(ctx)
	}
	if err != nil {
		// produce task for worker
		d.withdrawalProducer.Produce(ctx, withdrawId)
		return nil
	}

	// construct ledger entry
	ledger := &model.Ledger{
		BalanceId:     balance.Id,
		Type:          model.LedgerTypeOut,
		Amount:        withdrawal.Amount,
		BalanceBefore: balance.Amount,
	}

	if balance.Amount > withdrawal.Amount {
		// deduct balance
		balance.Amount -= withdrawal.Amount

		// assign new balance
		ledger.BalanceAfter = balance.Amount

		// update balance
		_, err = d.balanceRepo.Update(ctx, balance)
		if err != nil {
			return err
		}

		withdrawal.Status = model.StatusCompleted
	} else {
		withdrawal.Status = model.StatusFailed
	}

	// update withdrawal
	_, err = d.withdrawalRepo.Update(ctx, withdrawal)
	if err != nil {
		return err
	}

	if withdrawal.Status == model.StatusCompleted {
		// create ledger entry when withdrawal is ok
		ledger, err = d.ledgerRepo.Create(ctx, ledger)
		if err != nil {
			return err
		}
	}

	return nil
}
