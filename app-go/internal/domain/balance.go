package domain

import (
	"app/internal/model"
	"app/internal/repository"
	"app/internal/worker"
	"context"

	"github.com/segmentio/ksuid"
)

type IBalance interface {
	Open(ctx context.Context) (*model.Balance, error)
	Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error)
	Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error)
	Withdraw(ctx context.Context, in *model.WithdrawParam) (*model.Withdrawal, error)
	Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error)
}

type balance struct {
	balanceRepository    repository.IBalance
	depositRepository    repository.IDeposit
	withdrawalRepository repository.IWithdrawal
	transferRepository   repository.ITransfer
	withdrawalProducer   worker.IWithdrawalProducer
}

func NewBalance(
	balanceRepository repository.IBalance,
	depositRepository repository.IDeposit,
	withdrawalRepository repository.IWithdrawal,
	transferRepository repository.ITransfer,
	withdrawalProducer worker.IWithdrawalProducer,
) *balance {
	return &balance{
		balanceRepository:    balanceRepository,
		depositRepository:    depositRepository,
		withdrawalRepository: withdrawalRepository,
		transferRepository:   transferRepository,
		withdrawalProducer:   withdrawalProducer,
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

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	return d.transferRepository.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
}
