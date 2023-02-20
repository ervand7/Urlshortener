package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// UserURLsDelete /shortener.Shortener/UserURLsDelete.
func (s *ShortenerServer) UserURLsDelete(
	ctx context.Context, r *pb.UserURLsDeleteRequest,
) (*pb.UserURLsDeleteResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	go func() {
		s.Storage.DeleteUserURLs(r.URLs)
	}()
	response := pb.UserURLsDeleteResponse{Status: int32(codes.OK)}
	return &response, nil
}
