package domain

import (
	"app/internal/model"
	"app/internal/repository"
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
}

func NewBalance(
	balanceRepository repository.IBalance,
	depositRepository repository.IDeposit,
	withdrawalRepository repository.IWithdrawal,
	transferRepository repository.ITransfer,
) *balance {
	return &balance{
		balanceRepository:    balanceRepository,
		depositRepository:    depositRepository,
		withdrawalRepository: withdrawalRepository,
		transferRepository:   transferRepository,
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
	return d.withdrawalRepository.Create(ctx, &model.Withdrawal{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
}

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	return d.transferRepository.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
}
