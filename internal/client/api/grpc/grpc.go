package grpc

import (
	"context"
	"os"
	"time"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/protobuf/types/known/emptypb"
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

// AddCredentials - метод для создания учетных данных.
func (c *Client) AddCredentials(cred model.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	credPb := &pb.Credentials{}
	credPb.SetLogin(cred.Login)
	credPb.SetPassword(cred.Password)
	credPb.SetInfo(cred.Info)

	req := &pb.AddCredentialsRequest{}
	req.SetCredentials(credPb)

	_, err := c.grpcClient.AddCredentials(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetCredentials - метод для получения учетных данных.
func (c *Client) GetCredentials(alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetCredentialsRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.GetCredentials(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// DelCredentials - метод для удаления учетных данных.
func (c *Client) DelCredentials(alias []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.DelCredentialsRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.DelCredentials(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AddPaymentCard(card model.PaymentCard) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cardPb := &pb.PaymentCard{}
	cardPb.SetNumber(card.Number)
	cardPb.SetExpMonth(card.ExpMonth)
	cardPb.SetExpYear(card.ExpYear)
	cardPb.SetCvv(card.CVV)
	cardPb.SetInfo(card.Info)

	req := &pb.AddPaymentCardRequest{}
	req.SetCard(cardPb)

	_, err := c.grpcClient.AddPaymentCard(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetPaymentCard - метод для получения данных по карте.
func (c *Client) GetPaymentCard(alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetPaymentCardRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.GetPaymentCard(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// DelPaymentCard - метод для удаления данных по карте.
func (c *Client) DelPaymentCard(alias []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.DelPaymentCardRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.DelPaymentCard(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// AddBinary - метод для добавления бинарных данных.
func (c *Client) AddBinary(bin model.Binary) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	binPb := &pb.Binary{}
	binPb.SetData(bin.Data)
	binPb.SetInfo(bin.Info)

	req := &pb.AddBinaryRequest{}
	req.SetBin(binPb)

	_, err := c.grpcClient.AddBinary(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetBinary - метод для получения бинарных данных.
func (c *Client) GetBinary(alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetBinaryRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.GetBinary(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// DelBinary - метод для удаления бинарных данных.
func (c *Client) DelBinary(alias []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.DelBinaryRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.DelBinary(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// AddText - метод для добавления текстовых данных.
func (c *Client) AddText(text model.Text) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	textPb := &pb.Text{}
	textPb.SetData(text.Data)
	textPb.SetInfo(text.Info)

	req := &pb.AddTextRequest{}
	req.SetText(textPb)

	_, err := c.grpcClient.AddText(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// GetText - метод для получения текстовых данных.
func (c *Client) GetText(alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetTextRequest{}
	req.SetAlias(alias)

	_, err := c.grpcClient.GetText(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

// DelText - метод для удаления текстовых данных.
func (c *Client) DelText(alias []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.DelTextRequest{}
	req.SetAlias(alias)

	// var response *emptypb.Empty

	_, err := c.grpcClient.DelText(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListUserData(userId model.UserID) ([]model.Credentials, []model.PaymentCard, []model.Binary, []model.Text, err) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var req emptypb.Empty

	_, err := c.grpcClient.ListUserData(ctx, &req)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return nil, nil, nil, nil, nil
}
