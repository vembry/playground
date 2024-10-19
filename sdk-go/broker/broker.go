package broker

import (
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkpb "sdk/broker/pb"
)

type broker struct {
	conn   *grpc.ClientConn
	client sdkpb.BrokerClient
}

func New(addr string) (*broker, func()) {
	// setup grpc connection
	conn, err := grpc.NewClient(
		"localhost:4000", // Use the server's address and port
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := sdkpb.NewBrokerClient(conn)

	return &broker{
			conn:   conn,
			client: client,
		}, func() {
			conn.Close()
		}
}

func (b *broker) GRPC() sdkpb.BrokerClient {
	return b.client
}
