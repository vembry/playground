package domain

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessDeposit(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get withdraw
	deposit, err := d.depositRepository.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	balance, err := d.Get(ctx, deposit.BalanceId)
	if err != nil {
		return nil
	}

	// add balance
	balance.Amount += deposit.Amount

	// update balance
	_, err = d.balanceRepository.Update(ctx, balance)
	if err != nil {
		return err
	}

	deposit.Status = model.StatusCompleted

	// update withdrawal
	_, err = d.depositRepository.Update(ctx, deposit)
	if err != nil {
		return err
	}

	return nil
}
