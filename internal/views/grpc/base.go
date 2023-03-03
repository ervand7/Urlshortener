package grpc

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc/metadata"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// ShortenerServer struct for drpc server
type ShortenerServer struct {
	pb.UnimplementedShortenerServer
	Storage models.Storage
}

// GetUserIDFromMetadata gets UserID from grpc metadata
func (s ShortenerServer) GetUserIDFromMetadata(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("userid")
		if len(values) > 0 {
			decoded, err := hex.DecodeString(values[0])
			if err != nil {
				logger.Logger.Error(err.Error())
				return ""
			}
			return string(decoded)
		}
	}
	return ""
}
