package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ervand7/urlshortener/internal/logger"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// UserURLs /shortener.Shortener/UserURLs.
func (s *ShortenerServer) UserURLs(
	ctx context.Context, r *pb.UserURLsRequest,
) (*pb.UserURLsResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	userURLs, err := s.Storage.GetUserURLs(ctx, userID)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	response := pb.UserURLsResponse{}
	for _, item := range userURLs {
		response.Items = append(
			response.Items,
			&pb.UserURLsResponseItem{
				OriginalURL: item["original_url"],
				ShortURL:    item["short_url"],
			},
		)
	}

	return &response, nil
}
