package asynq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/segmentio/ksuid"
)

type withdrawal struct {
	client        *asynq.Client
	balanceDomain IWithdrawProcessor
}

func NewWithdrawal(client *asynq.Client) *withdrawal {
	return &withdrawal{
		client: client,
	}
}

func (w *withdrawal) Path() string {
	return "withdrawal"
}

func (w *withdrawal) Priority() int {
	return 1
}

type IWithdrawProcessor interface {
	ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (w *withdrawal) InjectDeps(balanceDomain IWithdrawProcessor) {
	w.balanceDomain = balanceDomain
}

func (w *withdrawal) Produce(ctx context.Context, withdrawalId ksuid.KSUID) error {
	payload, _ := withdrawalId.MarshalText()
	task := asynq.NewTask(
		w.Path(), payload,
		asynq.Queue(w.Path()),
	)

	_, err := w.client.EnqueueContext(ctx, task)
	return err
}

func (w *withdrawal) Handle(ctx context.Context, task *asynq.Task) error {
	id, _ := ksuid.Parse(string(task.Payload()))
	return w.balanceDomain.ProcessWithdraw(ctx, id)
}
