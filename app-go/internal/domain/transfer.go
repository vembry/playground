package domain

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessTransfer(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get withdraw
	transfer, err := d.transferRepo.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	balanceFrom, err := d.Get(ctx, transfer.BalanceIdFrom)
	if err != nil {
		return nil
	}

	ledgerFrom, ledgerTo := &model.Ledger{
		BalanceId:     transfer.BalanceIdFrom,
		Type:          model.LedgerTypeOut,
		Amount:        transfer.Amount,
		BalanceBefore: balanceFrom.Amount,
	}, &model.Ledger{
		BalanceId: transfer.BalanceIdTo,
		Type:      model.LedgerTypeIn,
		Amount:    transfer.Amount,
	}

	if balanceFrom.Amount > transfer.Amount {
		balanceTo, err := d.Get(ctx, transfer.BalanceIdTo)
		if err != nil {
			return nil
		}

		ledgerTo.BalanceBefore = balanceTo.Amount

		// deduct balance-from
		balanceFrom.Amount -= transfer.Amount
		// add balance-to
		balanceTo.Amount += transfer.Amount

		ledgerFrom.BalanceAfter = balanceFrom.Amount
		ledgerTo.BalanceAfter = balanceTo.Amount

		// update balance
		_, err = d.balanceRepo.Update(ctx, balanceFrom)
		if err != nil {
			return err
		}
		// update balance
		_, err = d.balanceRepo.Update(ctx, balanceTo)
		if err != nil {
			return err
		}

		transfer.Status = model.StatusCompleted
	} else {
		transfer.Status = model.StatusFailed
	}

	// update withdrawal
	_, err = d.transferRepo.Update(ctx, transfer)
	if err != nil {
		return err
	}

	if transfer.Status == model.StatusCompleted {
		// create ledger entry when transfer ok
		_, err = d.ledgerRepo.Create(ctx, ledgerFrom)
		if err != nil {
			return err
		}
		_, err = d.ledgerRepo.Create(ctx, ledgerTo)
		if err != nil {
			return err
		}
	}

	return nil
}
