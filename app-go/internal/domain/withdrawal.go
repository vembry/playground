package domain

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error {
	// get withdraw
	withdrawal, err := d.withdrawalRepository.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	balance, err := d.Get(ctx, withdrawal.BalanceId)
	if err != nil {
		return nil
	}

	if balance.Amount > withdrawal.Amount {
		// deduct balance
		balance.Amount -= withdrawal.Amount

		// update balance
		_, err = d.balanceRepository.Update(ctx, balance)
		if err != nil {
			return err
		}

		withdrawal.Status = model.StatusCompleted
	} else {
		withdrawal.Status = model.StatusFailed
	}

	// update withdrawal
	_, err = d.withdrawalRepository.Update(ctx, withdrawal)
	if err != nil {
		return err
	}

	return nil
}
