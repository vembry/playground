package domain

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessDeposit(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get deposit data
	deposit, err := d.depositRepo.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	// get balance
	balance, err := d.Get(ctx, deposit.BalanceId)
	if err != nil {
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
