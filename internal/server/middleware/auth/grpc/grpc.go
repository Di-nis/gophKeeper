// Package grpc реализует аутентификацию пользователя в gRPC сервере.
package grpc

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor - аутентификация пользователя.
func Interceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var (
			userID model.UserID
			token  string
			err    error
		)

		// TODO: найти куки
		mdReq, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := mdReq.Get(auth.HeaderAuthorization)
			if len(values) > 0 {
				token = values[0]
			}
		}

		if token == "" {
			return nil, status.Error(codes.Internal, "internal error")
		}

		claims, isTokenValid := auth.GetClaims(token, jwtSecret)
		if !isTokenValid {
			return nil, status.Error(codes.Unauthenticated, "token not valid")
		}
		userID = claims.UserID

		mdOut := metadata.Pairs(
			auth.HeaderAuthorization, token,
		)

		err = grpc.SendHeader(ctx, mdOut)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}

		ctx = context.WithValue(ctx, auth.Key, userID)
		return handler(ctx, req)
	}
}
