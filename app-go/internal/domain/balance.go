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
	lockerRepo         repository.ILocker
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
	lockerRepo repository.ILocker,
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
		lockerRepo:         lockerRepo,
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

func (d *balance) GetLock(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, func(context.Context), error) {
	var err error

	unlocker, err := d.lockerRepo.AcquireLock(ctx, balanceId.String())
	defer func() {
		if err != nil && unlocker != nil {
			unlocker(ctx)
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	balance, err := d.Get(ctx, balanceId)
	if err != nil {
		return nil, nil, err
	}

	return balance, func(_ctx context.Context) {
		unlocker(_ctx)
	}, nil
}

func (d *balance) Deposit(ctx context.Context, in *model.DepositParam) (*model.Deposit, error) {
	deposit, err := d.depositRepo.Create(ctx, &model.Deposit{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// produce task for worker
	err = d.depositProducer.Produce(ctx, deposit.Id)
	if err != nil {
		log.Printf("error on producing deposit task. depositId=%s. err=%v", deposit.Id, err)
	}

	return deposit, nil
}

func (d *balance) Withdraw(ctx context.Context, in *model.WithdrawParam) (*model.Withdrawal, error) {
	// create entry
	withdrawal, err := d.withdrawalRepo.Create(ctx, &model.Withdrawal{BalanceId: in.BalanceId, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// produce task for worker
	err = d.withdrawalProducer.Produce(ctx, withdrawal.Id)
	if err != nil {
		log.Printf("error on producing withdrawal task. withdrawalId=%s. err=%v", withdrawal.Id, err)
	}

	return withdrawal, nil
}

func (d *balance) Transfer(ctx context.Context, in *model.TransferParam) (*model.Transfer, error) {
	transfer, err := d.transferRepo.Create(ctx, &model.Transfer{BalanceIdFrom: in.BalanceIdFrom, BalanceIdTo: in.BalanceIdTo, Amount: in.Amount, Status: model.StatusPending})
	if err != nil {
		return nil, err
	}

	// produce task for worker
	err = d.transferProducer.Produce(ctx, transfer.Id)
	if err != nil {
		log.Printf("error on producing transfer task. transferId=%s. err=%v", transfer.Id, err)
	}

	return transfer, nil
}
