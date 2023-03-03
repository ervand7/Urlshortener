package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// ShortenBatch /shortener.Shortener/ShortenBatch.
func (s *ShortenerServer) ShortenBatch(
	ctx context.Context, r *pb.ShortenBatchRequest,
) (*pb.ShortenBatchResponse, error) {
	userID := s.GetUserIDFromMetadata(ctx)
	if userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "%s", "user ID is empty")
	}

	var dbEntries []models.Entry
	response := pb.ShortenBatchResponse{}
	for _, val := range r.Items {
		short := algorithms.GenerateShortURL()
		response.Items = append(
			response.Items, &pb.ShortenBatchResponseItem{
				CorrelationId: val.CorrelationId,
				ShortURL:      short,
			},
		)

		entry := models.Entry{UserID: userID, Short: short, Origin: val.OriginalURL}
		dbEntries = append(dbEntries, entry)
	}

	if err := s.Storage.SetMany(ctx, dbEntries); err != nil {
		logger.Logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}
	response.Status = int32(codes.OK)
	return &response, nil
}
