package service

import (
	"app/internal/model"
	"context"

	"github.com/segmentio/ksuid"
)

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	transfer, err := d.transferRepo.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// produce task for worker
	d.transferProducer.Produce(ctx, transfer.Id)

	return transfer, nil
}

func (d *balance) ProcessTransfer(ctx context.Context, transferId ksuid.KSUID) error {
	// get withdraw
	transfer, err := d.transferRepo.Get(ctx, transferId)
	if err != nil {
		return err
	}

	// validate state
	if transfer == nil || transfer.Status != model.StatusPending {
		return nil
	}

	// get balance-from lock
	balanceFrom, unlockerFrom, err := d.GetLock(ctx, transfer.BalanceIdFrom)
	if unlockerFrom != nil {
		defer unlockerFrom(ctx)
	}
	if err != nil {
		// produce task for worker
		d.transferProducer.Produce(ctx, transferId)
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
		// get balance-to lock
		balanceTo, unlockerTo, err := d.GetLock(ctx, transfer.BalanceIdTo)
		if unlockerTo != nil {
			defer unlockerTo(ctx)
		}
		if err != nil {
			// produce task for worker
			d.transferProducer.Produce(ctx, transferId)
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
