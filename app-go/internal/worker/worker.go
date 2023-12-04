package worker

import (
	"context"

	"github.com/segmentio/ksuid"
)

type IWithdrawalProducer interface {
	Produce(ctx context.Context, withdrawalId ksuid.KSUID) error
}
