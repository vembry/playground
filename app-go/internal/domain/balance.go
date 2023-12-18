package domain

import (
	"app/internal/model"
	"app/internal/repository"
	"app/internal/worker"
	"context"
	"log"

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
	balanceRepo        repository.IBalance
	depositRepo        repository.IDeposit
	withdrawalRepo     repository.IWithdrawal
	transferRepo       repository.ITransfer
	depositProducer    worker.IDepositProducer
	withdrawalProducer worker.IWithdrawalProducer
	transferProducer   worker.ITransferProducer
	ledgerRepo         repository.ILedger
}

func NewBalance(
	balanceRepo repository.IBalance,
	depositRepo repository.IDeposit,
	withdrawalRepo repository.IWithdrawal,
	transferRepo repository.ITransfer,
	depositProducer worker.IDepositProducer,
	withdrawalProducer worker.IWithdrawalProducer,
	transferProducer worker.ITransferProducer,
	ledgerRepo repository.ILedger,
) *balance {
	return &balance{
		balanceRepo:        balanceRepo,
		depositRepo:        depositRepo,
		withdrawalRepo:     withdrawalRepo,
		transferRepo:       transferRepo,
		depositProducer:    depositProducer,
		withdrawalProducer: withdrawalProducer,
		transferProducer:   transferProducer,
		ledgerRepo:         ledgerRepo,
	}
}

func (d *balance) Open(ctx context.Context) (*model.Balance, error) {
	return d.balanceRepo.Create(ctx, &model.Balance{
		Amount: float64(0),
	})
}

func (d *balance) Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error) {
	return d.balanceRepo.Get(ctx, balanceId)
}

func (d *balance) Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error) {
	deposit, err := d.depositRepo.Create(ctx, &model.Deposit{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// publish for worker
	err = d.depositProducer.Produce(ctx, deposit.Id)
	if err != nil {
		log.Printf("error on producing deposit task. err=%v", err)
	}

	return deposit, nil
}

func (d *balance) Withdraw(ctx context.Context, in *model.WithdrawParam) (*model.Withdrawal, error) {
	// create entry
	withdrawal, err := d.withdrawalRepo.Create(ctx, &model.Withdrawal{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// publish for worker
	err = d.withdrawalProducer.Produce(ctx, withdrawal.Id)
	if err != nil {
		log.Printf("error on producing withdrawal task. err=%v", err)
	}

	return withdrawal, nil
}

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	transfer, err := d.transferRepo.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// publish for worker
	err = d.transferProducer.Produce(ctx, transfer.Id)
	if err != nil {
		log.Printf("error on producing transfer task. err=%v", err)
	}

	return transfer, nil
}
