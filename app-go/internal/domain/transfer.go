package domain

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessTransfer(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get withdraw
	transfer, err := d.transferRepository.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	balanceFrom, err := d.Get(ctx, transfer.BalanceIdFrom)
	if err != nil {
		return nil
	}

	if balanceFrom.Amount > transfer.Amount {
		balanceTo, err := d.Get(ctx, transfer.BalanceIdTo)
		if err != nil {
			return nil
		}

		// deduct balance-from
		balanceFrom.Amount -= transfer.Amount
		// add balance-to
		balanceTo.Amount -= transfer.Amount

		// update balance
		_, err = d.balanceRepository.Update(ctx, balanceFrom)
		if err != nil {
			return err
		}
		// update balance
		_, err = d.balanceRepository.Update(ctx, balanceTo)
		if err != nil {
			return err
		}

		transfer.Status = model.StatusCompleted
	} else {
		transfer.Status = model.StatusFailed
	}

	// update withdrawal
	_, err = d.transferRepository.Update(ctx, transfer)
	if err != nil {
		return err
	}

	return nil
}
