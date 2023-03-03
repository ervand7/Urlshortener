package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/logger"
	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
	g "github.com/ervand7/urlshortener/internal/views/grpc"
)

// GetServer creates new grpc.Server
func GetServer() *grpc.Server {
	s := grpc.NewServer()
	return s
}

// Run launches grpc.Server
func Run(server *grpc.Server) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterShortenerServer(
		server,
		&g.ShortenerServer{
			Storage: storage.GetStorage(),
		},
	)

	logger.Logger.Info("===== RPC server started =====")
	if err = server.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
