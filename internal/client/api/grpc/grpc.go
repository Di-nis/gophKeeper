// Package grpc - реализация api gRPC-клиента.
package grpc

import (
	"context"
	"os"
	"time"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/middleware/auth"
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
func New(c *config.Config) (*Client, error) {
	conn, err := grpc.NewClient(
		c.ServerAddressGRPC,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(auth.Interceptor(c.TokenStorage)),
	)
	if err != nil {
		os.Exit(1)
	}
	// defer conn.Close()

	return &Client{
		grpcClient: pb.NewGophKeeperServiceClient(conn),
		timeout:    Timeout,
	}, nil
}

// AddCredentials - метод для создания учетных данных.
func (c *Client) AddCredentials(ctx context.Context, cred model.Credentials) (string, error) {
	// ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	// defer cancel()

	credPb := &pb.Credentials{}
	credPb.SetLogin(cred.Login)
	credPb.SetPassword(cred.Password)
	credPb.SetInfo(cred.Info)

	req := &pb.AddCredentialsRequest{}
	req.SetCredentials(credPb)

	response, err := c.grpcClient.AddCredentials(ctx, req)
	if err != nil {
		return "", err
	}
	return response.GetAlias(), nil
}

// GetCredentials - метод для получения учетных данных.
func (c *Client) GetCredentials(alias string) (model.Credentials, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetCredentialsRequest{}
	req.SetAlias(alias)

	cred := model.Credentials{}
	response, err := c.grpcClient.GetCredentials(ctx, req)
	if err != nil {
		return cred, err
	}

	cred.Login = response.GetCredentials().GetLogin()
	cred.Password = response.GetCredentials().GetPassword()
	cred.Info = response.GetCredentials().GetInfo()

	return cred, nil
}

// DelCredentials - метод для удаления учетных данных.
func (c *Client) DelCredentials(alias string) error {
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

func (c *Client) AddPaymentCard(card model.PaymentCard) (string, error) {
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

	response, err := c.grpcClient.AddPaymentCard(ctx, req)
	if err != nil {
		return "", err
	}
	return response.GetAlias(), nil
}

// GetPaymentCard - метод для получения данных по карте.
func (c *Client) GetPaymentCard(alias string) (model.PaymentCard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetPaymentCardRequest{}
	req.SetAlias(alias)

	card := model.PaymentCard{}
	response, err := c.grpcClient.GetPaymentCard(ctx, req)
	if err != nil {
		return card, err
	}

	card.Number = response.GetCard().GetNumber()
	card.ExpMonth = response.GetCard().GetExpMonth()
	card.ExpYear = response.GetCard().GetExpYear()
	card.CVV = response.GetCard().GetCvv()
	card.Holder = response.GetCard().GetHolder()
	card.Info = response.GetCard().GetInfo()

	return card, nil
}

// DelPaymentCard - метод для удаления данных по карте.
func (c *Client) DelPaymentCard(alias string) error {
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
func (c *Client) AddBinary(bin model.Binary) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	binPb := &pb.Binary{}
	binPb.SetData(bin.Data)
	binPb.SetInfo(bin.Info)

	req := &pb.AddBinaryRequest{}
	req.SetBin(binPb)

	response, err := c.grpcClient.AddBinary(ctx, req)
	if err != nil {
		return "", err
	}
	return response.GetAlias(), nil
}

// GetBinary - метод для получения бинарных данных.
func (c *Client) GetBinary(alias string) (model.Binary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetBinaryRequest{}
	req.SetAlias(alias)

	bin := model.Binary{}
	response, err := c.grpcClient.GetBinary(ctx, req)
	if err != nil {
		return bin, err
	}

	bin.Data = response.GetBin().GetData()
	bin.Info = response.GetBin().GetInfo()

	return bin, nil
}

// DelBinary - метод для удаления бинарных данных.
func (c *Client) DelBinary(alias string) error {
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
func (c *Client) AddText(text model.Text) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	textPb := &pb.Text{}
	textPb.SetData(text.Data)
	textPb.SetInfo(text.Info)

	req := &pb.AddTextRequest{}
	req.SetText(textPb)

	response, err := c.grpcClient.AddText(ctx, req)
	if err != nil {
		return "", err
	}
	return response.GetAlias(), nil
}

// GetText - метод для получения текстовых данных.
func (c *Client) GetText(alias string) (model.Text, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.GetTextRequest{}
	req.SetAlias(alias)

	text := model.Text{}
	response, err := c.grpcClient.GetText(ctx, req)
	if err != nil {
		return text, err
	}

	text.Data = response.GetText().GetData()
	text.Info = response.GetText().GetInfo()
	return text, nil
}

// DelText - метод для удаления текстовых данных.
func (c *Client) DelText(alias string) error {
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

// Sync - синхронизация данных.
func (c *Client) Sync() (model.Common, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var req emptypb.Empty

	response, err := c.grpcClient.ListUserData(ctx, &req)
	if err != nil {
		return model.Common{}, err
	}

	credsPb := response.GetCredentials()
	cardsPb := response.GetPaymentCard()
	binsPb := response.GetBinary()
	textsPb := response.GetText()

	creds := convertCreds(credsPb)
	cards := convertCards(cardsPb)
	bins := convertBinary(binsPb)
	texts := convertText(textsPb)

	return model.Common{
		Credentials: creds,
		PaymentCard: cards,
		Binary:      bins,
		Text:        texts,
	}, nil
}
