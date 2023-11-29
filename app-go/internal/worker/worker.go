package worker

import "context"

type IWorker interface {
	Name() string
	Start()
}

type IProducer interface {
	IWorker
	Produce(ctx context.Context)
}
