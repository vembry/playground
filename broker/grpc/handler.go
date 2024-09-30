package grpc

import (
	"broker/model"
	"context"
	"fmt"
	"sdk/pb"

	"github.com/segmentio/ksuid"
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

// GetQueue gets all queues data
func (h *handler) GetQueue(ctx context.Context, req *pb.GetQueueRequest) (*pb.GetQueueResponse, error) {
	res := h.queue.Get()

	activeQueues := map[string]*pb.ActiveQueue{}
	for _, activeQueue := range res.ActiveQueue {
		activeQueues[activeQueue.Id.String()] = &pb.ActiveQueue{
			Id:         activeQueue.Id.String(),
			QueueName:  activeQueue.QueueName,
			PollExpiry: activeQueue.PollExpiry.String(),
			Payload:    activeQueue.Payload,
		}
	}

	idleQueueMap := map[string]*pb.IdleQueue{}
	for key, val := range res.IdleQueue {
		idleQueueMap[key] = &pb.IdleQueue{
			Items: val.Items,
		}
	}

	return &pb.GetQueueResponse{
		Message: "ok",
		Data: &pb.QueueData{
			ActiveQueue: activeQueues,
			IdleQueue:   idleQueueMap,
		},
	}, nil
}

// Enqueue enqueues entry to the queue
func (h *handler) Enqueue(ctx context.Context, req *pb.EnqueueRequest) (*pb.EnqueueResponse, error) {
	err := h.queue.Enqueue(model.EnqueuePayload{
		Name:    req.GetQueueName(),
		Payload: req.GetPayload(),
	})
	if err != nil {
		return nil, fmt.Errorf("error on enqueue")
	}
	return &pb.EnqueueResponse{
		Message: "ok",
	}, nil
}

// Poll retrieves selected queue's entry
func (h *handler) Poll(ctx context.Context, req *pb.PollRequest) (*pb.PollResponse, error) {
	queue, err := h.queue.Poll(req.GetQueueName())
	if err != nil {
		return &pb.PollResponse{
			Message: err.Error(),
		}, fmt.Errorf("error to poll")
	}

	// when no entry can be polled from queue
	// then return nothing
	if queue == nil {
		return &pb.PollResponse{
			Message: "no queue",
			Data:    nil,
		}, nil
	}

	return &pb.PollResponse{
		Message: "ok",
		Data: &pb.ActiveQueue{
			Id:         queue.Id.String(),
			QueueName:  queue.QueueName,
			PollExpiry: queue.PollExpiry.String(),
			Payload:    queue.Payload,
		},
	}, nil
}

// CompletePoll acks polled queue entry
func (h *handler) CompletePoll(ctx context.Context, req *pb.CompletePollRequest) (*pb.CompletePollResponse, error) {
	queueId, err := ksuid.Parse(req.GetQueueId())
	if err != nil {
		return nil, fmt.Errorf("invalid queue id")
	}

	err = h.queue.CompletePoll(queueId)
	if err != nil {
		return nil, fmt.Errorf("failed to complete-poll")
	}

	return &pb.CompletePollResponse{
		Message: "ok",
	}, nil
}
