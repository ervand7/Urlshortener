package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
	"github.com/ervand7/urlshortener/internal/views"
)

// GetStats /shortener.Shortener/GetStats.
func (s *ShortenerServer) GetStats(
	ctx context.Context, r *pb.GetStatsRequest,
) (*pb.GetStatsResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	stats, err := views.GetStats(ctx, s.Storage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	response := pb.GetStatsResponse{
		Urls:  int64(stats.NumberOfURLs),
		Users: int64(stats.NumberOfUsers),
	}
	return &response, nil
}
