// Package auth - middleware для gRPC сервера.
package auth

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Interceptor - аутентификация пользователя.
func Interceptor(path string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		token, err := ReadFileAsString(path)
		if err != nil {
			return err
		}

		ctx = metadata.AppendToOutgoingContext(ctx, auth.HeaderAuthorization, token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
