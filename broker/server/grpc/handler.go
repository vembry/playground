package grpc

import (
	"broker/model"
	"broker/server"
	"context"
	"fmt"
	"sdk/broker/pb"

	"github.com/segmentio/ksuid"
)

type handler struct {
	pb.UnimplementedBrokerServer
	broker server.IBroker
}

func NewHandler(broker server.IBroker) *handler {
	return &handler{
		broker: broker,
	}
}

// GetQueue gets all queues data
func (h *handler) GetQueue(ctx context.Context, req *pb.GetQueueRequest) (*pb.GetQueueResponse, error) {
	res := h.broker.Get()

	return &pb.GetQueueResponse{
		Message: "ok",
		Data: &pb.QueueData{
			IdleQueueCount:   res.IdleQueueCount,
			ActiveQueueCount: res.ActiveQueueCount,
		},
	}, nil
}

// Enqueue enqueues entry to the queue
func (h *handler) Enqueue(ctx context.Context, req *pb.EnqueueRequest) (*pb.EnqueueResponse, error) {
	err := h.broker.Enqueue(model.EnqueuePayload{
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
	queue, err := h.broker.Poll(req.GetQueueName())
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
			Queue: &pb.Queue{
				Payload: queue.Queue.Payload,
			},
		},
	}, nil
}

// CompletePoll acks polled queue entry
func (h *handler) CompletePoll(ctx context.Context, req *pb.CompletePollRequest) (*pb.CompletePollResponse, error) {
	queueId, err := ksuid.Parse(req.GetQueueId())
	if err != nil {
		return nil, fmt.Errorf("invalid queue id")
	}

	err = h.broker.CompletePoll(queueId)
	if err != nil {
		return nil, fmt.Errorf("failed to complete-poll")
	}

	return &pb.CompletePollResponse{
		Message: "ok",
	}, nil
}
