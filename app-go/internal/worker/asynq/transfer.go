package asynq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/segmentio/ksuid"
)

type transfer struct {
	client            *asynq.Client
	transferProcessor ITransferProcessor
}

func NewTransfer(client *asynq.Client) *transfer {
	return &transfer{
		client: client,
	}
}

func (w *transfer) Path() string {
	return "transfer"
}

func (w *transfer) Priority() int {
	return 1
}

type ITransferProcessor interface {
	ProcessTransfer(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (w *transfer) InjectDeps(transferProcessor ITransferProcessor) {
	w.transferProcessor = transferProcessor
}

func (w *transfer) Produce(ctx context.Context, transferId ksuid.KSUID) error {
	payload, _ := transferId.MarshalText()
	task := asynq.NewTask(
		w.Path(), payload,
		asynq.Queue(w.Path()),
	)

	_, err := w.client.EnqueueContext(ctx, task)
	return err
}

func (w *transfer) Handle(ctx context.Context, task *asynq.Task) error {
	id, _ := ksuid.Parse(string(task.Payload()))
	return w.transferProcessor.ProcessTransfer(ctx, id)
}
