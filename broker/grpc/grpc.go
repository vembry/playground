package grpc

import (
	"broker/model"
	"log"
	"net"
	"sdk/broker/pb"

	"github.com/segmentio/ksuid"
	"google.golang.org/grpc"
)

type IQueue interface {
	Get() model.QueueData
	Enqueue(payload model.EnqueuePayload) error
	Poll(queueName string) (*model.ActiveQueue, error)
	CompletePoll(queueId ksuid.KSUID) error
}

type server struct {
	server *grpc.Server
}

func NewServer(queue IQueue) *server {
	grpcServer := grpc.NewServer()
	handler := NewHandler(queue)

	pb.RegisterBrokerServer(grpcServer, handler)

	return &server{
		server: grpcServer,
	}
}

func (g *server) Start() error {
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return g.server.Serve(lis)
}
func (g *server) Stop() {
	g.server.GracefulStop()
}
