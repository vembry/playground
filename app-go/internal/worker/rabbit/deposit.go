package rabbit

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
)

type deposit struct {
	ch               *amqp.Channel
	depositProcessor IDepositProcessor
}

func NewDeposit(conn *amqp.Connection) *deposit {
	// setup channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. err=%v", err)
	}

	return &deposit{
		ch: ch,
	}
}

func (t *deposit) Name() string {
	return "deposit"
}
func (t *deposit) Channel() *amqp.Channel {
	return t.ch
}

func (t *deposit) Produce(ctx context.Context, depositId ksuid.KSUID) error {
	body, _ := depositId.MarshalText()

	// enqueue task
	err := t.ch.PublishWithContext(ctx,
		"",       // exchange
		t.Name(), // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish a queue's task. err=%v", err)
	}
	return nil
}

type IDepositProcessor interface {
	ProcessDeposit(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (t *deposit) InjectDeps(depositProcessor IDepositProcessor) {
	t.depositProcessor = depositProcessor
}

func (t *deposit) Handle(ctx context.Context, task amqp.Delivery) error {
	id, _ := ksuid.Parse(string(task.Body))
	return t.depositProcessor.ProcessDeposit(ctx, id)
}
