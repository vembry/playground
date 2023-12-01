package kafka

import (
	"app/internal/app"
	"context"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/ksuid"
)

type withdrawal struct {
	reader        *kafkago.Reader
	writer        *kafkago.Writer
	balanceDomain IBalance
}

func NewWithdrawal(cfg *app.EnvConfig) *withdrawal {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    "withdrawal",
		GroupID:  "app-go",
		MaxBytes: 10e6, // 10MB
	})

	w := kafkago.NewWriter(kafkago.WriterConfig{
		Brokers:      []string{cfg.KafkaBroker},
		Topic:        "withdrawal",
		Balancer:     &kafkago.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	})

	return &withdrawal{
		reader: r,
		writer: w,
	}
}

func (w *withdrawal) Name() string {
	return "kafka.withdrawal"
}

func (w *withdrawal) Produce(ctx context.Context, withdrawalId ksuid.KSUID) error {
	return w.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(w.Name()),
		Value: []byte(withdrawalId.String()),
	})
}

func (w *withdrawal) Start() {
	log.Printf("starting %s", w.Name())
	for {
		ctx := context.Background()
		msg, err := w.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error on reading message. err=%v", err)
		}
		log.Printf("received message. value=%v", string(msg.Value))

		withdrawId, _ := ksuid.Parse(string(msg.Value))

		err = w.balanceDomain.ProcessWithdraw(ctx, withdrawId)
		if err != nil {
			log.Printf("error processing withdrawal. err=%v", err)
		}
	}
}

func (w *withdrawal) Stop() {
	log.Printf("stopping %s", w.Name())
	w.reader.Close()
}

type IBalance interface {
	ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (w *withdrawal) InjectDep(balanceDomain IBalance) {
	w.balanceDomain = balanceDomain
}
