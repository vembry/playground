package worker

import (
	"context"

	"github.com/segmentio/ksuid"
)

type IWithdrawalProducer interface {
	Produce(ctx context.Context, withdrawalId ksuid.KSUID) error
}

type IDepositProducer interface {
	Produce(ctx context.Context, depositId ksuid.KSUID) error
}

type ITransferProducer interface {
	Produce(ctx context.Context, transferId ksuid.KSUID) error
}
