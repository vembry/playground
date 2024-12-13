package balance

import (
	"app/internal/model"
	"app/internal/module"
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
	balanceRepo        repository.IBalance
	depositRepo        repository.IDeposit
	withdrawalRepo     repository.IWithdrawal
	transferRepo       repository.ITransfer
	depositProducer    worker.IDepositProducer
	withdrawalProducer worker.IWithdrawalProducer
	transferProducer   worker.ITransferProducer
	ledgerRepo         repository.ILedger
	locker             module.ILocker
}

func New(
	balanceRepo repository.IBalance,
	depositRepo repository.IDeposit,
	withdrawalRepo repository.IWithdrawal,
	transferRepo repository.ITransfer,
	depositProducer worker.IDepositProducer,
	withdrawalProducer worker.IWithdrawalProducer,
	transferProducer worker.ITransferProducer,
	ledgerRepo repository.ILedger,
	locker module.ILocker,
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
		locker:             locker,
	}
}

// Open creates new balance entry
func (d *balance) Open(ctx context.Context) (*model.Balance, error) {
	return d.balanceRepo.Create(ctx, &model.Balance{
		Amount: float64(0),
	})
}

// Get gets balance by id
func (d *balance) Get(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, error) {
	return d.balanceRepo.Get(ctx, balanceId)
}

// TODO: this need review :/
func (d *balance) GetLock(ctx context.Context, balanceId ksuid.KSUID) (*model.Balance, func(context.Context), error) {
	var err error

	unlocker, err := d.locker.Lock(ctx, balanceId.String())
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
