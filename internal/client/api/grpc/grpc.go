package grpc

import (
	"context"
	"os"
	"time"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Timeout time.Duration = 5 * time.Second
)

// Client - структура gRPC-клиента.
type Client struct {
	grpcClient pb.GophKeeperServiceClient
	timeout    time.Duration
}

// New - создание нового клиента.
func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// slog.Error("ошибка при установлении соединения с сервером", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	return &Client{
		grpcClient: pb.NewGophKeeperServiceClient(conn),
		timeout:    Timeout,
	}, nil
}

// Credentials - метод для создания учетных данных.
func (c *Client) Credentials(credentialsIn model.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	credentialsPb := &pb.AddCredentialsRequest{
		Credentials: &pb.Credentials{
			Login:    &credentialsIn.Login,
			Password: &credentialsIn.Password,
			Info:     &credentialsIn.Info,
		},
	}

	_, err := c.grpcClient.AddCredentials(ctx, credentialsPb)
	if err != nil {
		return err
	}
	return nil
}
