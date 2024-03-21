package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/danblok/cleanbg/internal/pb"
	"github.com/danblok/cleanbg/internal/types"
)

type cleanerClient struct {
	client pb.CleanerServiceClient
}

func Connect(target string) (types.Cleaner, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10<<20)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewCleanerServiceClient(conn)
	return &cleanerClient{
		client: client,
	}, nil
}

func (c *cleanerClient) Clean(ctx context.Context, image []byte) ([]byte, error) {
	res, err := c.client.Clean(ctx,
		&pb.CleanRequest{
			Data: image,
		},
	)

	return res.Data, err
}
