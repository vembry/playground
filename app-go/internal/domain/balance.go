package domain

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"

	"github.com/segmentio/ksuid"
)

type IBalance interface {
	Open(ctx context.Context) (*model.Balance, error)
	Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error)
	Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error)
	Withdraw(ctx context.Context, in *model.WithdrawParam) (*model.Withdrawal, error)
	Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error)
}

type IWithdrawalProducer interface {
	Produce(ctx context.Context, withdrawalId ksuid.KSUID) error
}

type balance struct {
	balanceRepository    repository.IBalance
	depositRepository    repository.IDeposit
	withdrawalRepository repository.IWithdrawal
	transferRepository   repository.ITransfer
	withdrawalProducer   IWithdrawalProducer
	mutex                IMutex
}

func NewBalance(
	balanceRepository repository.IBalance,
	depositRepository repository.IDeposit,
	withdrawalRepository repository.IWithdrawal,
	transferRepository repository.ITransfer,
	withdrawalProducer IWithdrawalProducer,
	mutex IMutex,
) *balance {
	return &balance{
		balanceRepository:    balanceRepository,
		depositRepository:    depositRepository,
		withdrawalRepository: withdrawalRepository,
		transferRepository:   transferRepository,
		withdrawalProducer:   withdrawalProducer,
		mutex:                mutex,
	}
}

func (d *balance) Open(ctx context.Context) (*model.Balance, error) {
	return d.balanceRepository.Create(ctx, &model.Balance{
		Amount: float64(0),
	})
}

func (d *balance) Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error) {
	return d.balanceRepository.Get(ctx, balanceId)
}

func (d *balance) GetLock(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, func(), error) {
	balance, err := d.Get(ctx, balanceId)
	if err != nil {
		return nil, nil, err
	}

	mutex, err := d.mutex.Acquire(ctx, fmt.Sprintf("balance.%s", balanceId))
	if err != nil {
		return nil, nil, err
	}

	return balance, func() {
		d.mutex.Delete(ctx, mutex)
	}, nil
}

func (d *balance) Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error) {
	return d.depositRepository.Create(ctx, &model.Deposit{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
}

func (d *balance) Withdraw(ctx context.Context, in *model.WithdrawParam) (*model.Withdrawal, error) {
	// create entry
	withdrawal, err := d.withdrawalRepository.Create(ctx, &model.Withdrawal{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// publish for worker
	err = d.withdrawalProducer.Produce(ctx, withdrawal.Id)
	if err != nil {
		return nil, err
	}

	return withdrawal, nil
}

func (d *balance) ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error {
	// // lock withdrawal
	// mutex, err := d.mutex.Acquire(ctx, fmt.Sprintf("withdraw.%s", withdrawId))
	// if err != nil {
	// 	return err
	// }
	// if mutex == nil {
	// 	return nil
	// }
	// defer d.mutex.Delete(ctx, mutex)

	// get withdraw
	withdrawal, err := d.withdrawalRepository.Get(ctx, withdrawId)
	if err != nil {
		return err
	}

	// // get balance with lock
	// balance, releaseBalance, err := d.GetLock(ctx, withdrawal.BalanceId)
	// if err != nil {
	// 	return nil
	// }
	// if balance == nil {
	// 	log.Printf("balance not found")
	// 	return nil
	// }
	// defer releaseBalance()

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

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	return d.transferRepository.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
}
