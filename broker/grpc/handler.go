package grpc

import (
	"context"
	"sdk/pb"
)

type handler struct {
	pb.UnimplementedBrokerServer
	queue IQueue
}

func NewHandler(queue IQueue) *handler {
	return &handler{
		queue: queue,
	}
}

func (ig *handler) GetQueue(ctx context.Context, req *pb.GetQueueRequest) (*pb.GetQueueResponse, error) {
	res := ig.queue.Get()

	activeQueues := map[string]*pb.ActiveQueue{}
	for _, activeQueue := range res.ActiveQueue {
		activeQueues[activeQueue.Id.String()] = &pb.ActiveQueue{
			Id:         activeQueue.Id.String(),
			QueueName:  activeQueue.QueueName,
			PollExpiry: activeQueue.PollExpiry.String(),
			Payload:    activeQueue.Payload,
		}
	}

	queueList := map[string]*pb.QueueList{}
	for key, val := range res.Queue {
		queueList[key] = &pb.QueueList{
			Items: val.Items,
		}
	}

	return &pb.GetQueueResponse{
		Message: "ok",
		Data: &pb.QueueData{
			ActiveQueue: activeQueues,
			Queue:       queueList,
		},
	}, nil
}
