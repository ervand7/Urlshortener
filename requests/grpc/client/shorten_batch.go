package main

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

func shortenBatch(conn *grpc.ClientConn) (*pb.ShortenBatchResponse, error) {
	encoded := hex.EncodeToString([]byte(userID))
	md := metadata.New(map[string]string{"userid": encoded})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewShortenerClient(conn)
	resp, err := client.ShortenBatch(
		ctx,
		&pb.ShortenBatchRequest{
			Items: []*pb.ShortenBatchRequestItem{
				{
					CorrelationId: "799eb140-c1f8-45c4-8268-656c5b20595b",
					OriginalURL:   randomString(60),
				},
				{
					CorrelationId: "199eb140-c1f8-45c4-8268-656c5b20595b",
					OriginalURL:   randomString(60),
				},
			},
		},
	)
	return resp, err
}
