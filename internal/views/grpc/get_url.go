package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	e "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// GetURL /shortener.Shortener/GetURL.
func (s *ShortenerServer) GetURL(
	ctx context.Context, r *pb.GetURLRequest,
) (*pb.GetURLResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	origin, err := s.Storage.Get(ctx, r.ShortURL)
	if err != nil {
		logger.Logger.Error(err.Error())
		if _, ok := err.(*e.URLNotActiveError); ok {
			return nil, status.Errorf(codes.ResourceExhausted, "%s", err.Error())
		}
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, status.Errorf(codes.Unknown, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	response := pb.GetURLResponse{
		OriginURL: origin,
		Status:    int32(codes.OK),
	}
	return &response, nil
}
