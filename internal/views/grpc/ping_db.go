package grpc

import (
	"context"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
	"github.com/ervand7/urlshortener/internal/logger"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

// PingDB "/shortener.Shortener/PingDB".
func (s *ShortenerServer) PingDB(
	_ context.Context, _ *pb.PingDBRequest,
) (*pb.PingDBResponse, error) {
	db := dbstorage.Database{}
	defer func() {
		if err := os.Unsetenv("DATABASE_DSN"); err != nil {
			logger.Logger.Error(err.Error())
		}
		if err := db.ConnClose(); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()

	err := db.ConnStart()
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	err = db.Conn.Ping()
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	} else {
		response := pb.PingDBResponse{Status: int32(codes.OK)}
		return &response, nil
	}
}
