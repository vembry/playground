package service

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

// Deposit creates deposit entry and queue it to the background
func (d *balance) Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error) {
	deposit, err := d.depositRepo.Create(ctx, &model.Deposit{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// produce task for worker
	d.depositProducer.Produce(ctx, deposit.Id)

	return deposit, nil
}

// ProcessDeposit is handler to actually process deposit
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
	_, err = d.ledgerRepo.Create(ctx, ledger)
	if err != nil {
		return err
	}

	return nil
}
