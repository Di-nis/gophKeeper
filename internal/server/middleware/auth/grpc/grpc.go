package grpc

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/server/middleware/auth"

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
			token, userID, sessionID string
			err                      error
		)

		mdReq, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := mdReq.Get(auth.HeaderAuthorization)
			if len(values) > 0 {
				token = values[0]
			}
		}

		if token == "" {
			userID = auth.GenerateUserID()
			sessionID = auth.GenerateSessionID()
			token, err = auth.BuildJWTString(jwtSecret, userID, sessionID)
			if err != nil {
				return nil, status.Error(codes.Internal, "internal error")
			}
		} else {
			claims, isTokenValid := auth.GetClaims(token, jwtSecret)
			if !isTokenValid {
				return nil, status.Error(codes.Unauthenticated, "token not valid")
			}
			userID = claims.UserID
			sessionID = claims.SID
			if sessionID == "" {
				return nil, status.Error(codes.Unauthenticated, "sessionID not valid")
			}
		}

		mdOut := metadata.Pairs(
			auth.HeaderAuthorization, token,
		)

		err = grpc.SendHeader(ctx, mdOut)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}

		ctx = context.WithValue(ctx, auth.UserIDKey, userID)
		return handler(ctx, req)
	}
}
