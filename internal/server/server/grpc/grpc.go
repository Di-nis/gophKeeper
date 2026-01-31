// package grpc - gRPC-сервер.
package grpc

import (
	"context"
	"net"
	"os"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/server/config"
	handler "github.com/Di-nis/gophKeeper/internal/server/delivery/grpc"
	auth "github.com/Di-nis/gophKeeper/internal/server/middleware/auth/grpc"

	// "github.com/Di-nis/gophKeeper/internal/server/middleware/compress"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"google.golang.org/grpc"
)

// Server - структура gRPC-сервера.
type Server struct {
	grpc    *grpc.Server
	lis     net.Listener
	handler *handler.Handler
}

// New - создание gRPC-сервера.
func New(config *config.Config, handler *handler.Handler) (*Server, error) {
	lis, err := net.Listen("tcp", config.ServerAddressGRPC)
	if err != nil {
		logger.Sugar.Errorf("failed initializing listener, error - %v", err)
		os.Exit(1)
	}

	// TODO: add compress
	// server := grpc.NewServer(grpc.UnaryInterceptor(auth.Interceptor(config.JWTSecret)), grpc.UnaryInterceptor(compress.Interceptor()))
	server := grpc.NewServer(grpc.UnaryInterceptor(auth.Interceptor(config.JWTSecret)))
	pb.RegisterGophKeeperServiceServer(server, handler)

	return &Server{
		grpc:    server,
		lis:     lis,
		handler: handler,
	}, nil
}

// Start - запуск gRPC-сервера.
func (s *Server) Start() error {
	return s.grpc.Serve(s.lis)
}

// Stop - остановка gRPC-сервера.
func (s *Server) Stop(ctx context.Context) error {
	s.grpc.GracefulStop()
	return nil
}
