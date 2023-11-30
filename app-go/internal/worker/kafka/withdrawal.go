package kafka

import (
	"context"
	"fmt"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

type withdrawal struct {
	reader *kafkago.Reader
}

func NewWithdrawal() *withdrawal {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{"localhost:9093"},
		Topic:    "withdrawal",
		GroupID:  "app-go",
		MaxBytes: 10e6, // 10MB
	})

	return &withdrawal{
		reader: r,
	}
}

func (w *withdrawal) Name() string {
	return "kafka.withdrawal"
}

func (w *withdrawal) Start() {
	log.Printf("starting %s", w.Name())
	for {
		msg, err := w.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			// break
		}

		fmt.Printf("Received message: %s\n", msg.Value)
	}
}

func (w *withdrawal) Stop() {
	log.Printf("stopping %s", w.Name())
	w.reader.Close()
}
