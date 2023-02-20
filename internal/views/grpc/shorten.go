package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	e "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// ShortenURL /shortener.Shortener/ShortenURL.
func (s *ShortenerServer) ShortenURL(
	ctx context.Context, r *pb.ShortenURLRequest,
) (*pb.ShortenURLResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	short := algorithms.GenerateShortURL()
	if err := s.Storage.Set(ctx, userID, short, r.OriginURL); err != nil {
		logger.Logger.Error(err.Error())
		if _, ok := err.(*e.ShortAlreadyExistsError); ok {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
		} else {
			return nil, status.Errorf(codes.Internal, "%s", err.Error())
		}
	}

	response := pb.ShortenURLResponse{
		ShortenURL: short,
		Status:     int32(codes.OK),
	}
	return &response, nil
}
