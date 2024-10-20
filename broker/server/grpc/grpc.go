package grpc

import (
	"broker/server"
	"log"
	"net"
	"sdk/broker/pb"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type grpcserver struct {
	server *grpc.Server
}

func New(broker server.IBroker) *grpcserver {
	// create grpc server
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	handler := NewHandler(broker)

	pb.RegisterBrokerServer(grpcServer, handler)

	return &grpcserver{
		server: grpcServer,
	}
}

func (g *grpcserver) Start() error {
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return g.server.Serve(lis)
}
func (g *grpcserver) Stop() {
	g.server.GracefulStop()
}
