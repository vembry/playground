package asynq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/segmentio/ksuid"
)

type deposit struct {
	client           *asynq.Client
	depositProcessor IDepositProcessor
}

func NewDeposit(client *asynq.Client) *deposit {
	return &deposit{
		client: client,
	}
}

func (w *deposit) Path() string {
	return "deposit"
}

func (w *deposit) Priority() int {
	return 1
}

type IDepositProcessor interface {
	ProcessDeposit(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (w *deposit) InjectDeps(depositProcessor IDepositProcessor) {
	w.depositProcessor = depositProcessor
}

func (w *deposit) Produce(ctx context.Context, depositId ksuid.KSUID) error {
	payload, _ := depositId.MarshalText()
	task := asynq.NewTask(
		w.Path(), payload,
		asynq.Queue(w.Path()),
	)

	_, err := w.client.EnqueueContext(ctx, task)
	return err
}

func (w *deposit) Handle(ctx context.Context, task *asynq.Task) error {
	id, _ := ksuid.Parse(string(task.Payload()))
	return w.depositProcessor.ProcessDeposit(ctx, id)
}
