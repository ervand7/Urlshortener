package main

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

func getStats(conn *grpc.ClientConn) (*pb.GetStatsResponse, error) {
	encoded := hex.EncodeToString([]byte(userID))
	md := metadata.New(map[string]string{"userid": encoded})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewShortenerClient(conn)

	resp, err := client.GetStats(ctx, &pb.GetStatsRequest{})
	return resp, err
}
