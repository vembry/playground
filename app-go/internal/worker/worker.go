package worker

type IWorker interface {
	Name() string
}

type IConsumer interface {
	IWorker
	Start()
	Stop()
}
