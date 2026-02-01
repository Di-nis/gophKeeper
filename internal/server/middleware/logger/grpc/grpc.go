// Package grpc реализует аутентификацию пользователя в gRPC сервере.
package grpc

import (
	"context"
	"time"

	"github.com/Di-nis/gophKeeper/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Interceptor - аутентификация пользователя.
func Interceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)

		logger.Sugar.Infoln(
			"method", info.FullMethod,
			"status", st.Code(),
			"duration", time.Since(start),
		)

		return resp, err
	}
}
