// Package grpc - реализация gRPC-сервера.
package grpc

import (
	"context"
	"errors"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/config"

	usecase "github.com/Di-nis/gophKeeper/internal/server/usecase/data"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Creater - интерфейс для создания данных в базе данных.
type Creater interface {
	Create(context.Context, any) error
}

// Getter - интерфейс, включащий методы по получению всех данных пользователя.
type Getter interface {
	Get(context.Context, any) error
	GetAll(context.Context, model.UserID) ([]model.Credentials, []model.PaymentCard, []model.Binary, []model.Text, error)
}

// Deleter - интерфейс для удаления данных из базы данных.
type Deleter interface {
	Delete(context.Context, []any) error
}

// Service - интерфейс для работы с данными.
type Service interface {
	Creater
	Getter
	Deleter
}

// Handler поддерживает все необходимые методы сервера.
type Handler struct {
	pb.GophKeeperServiceServer
	usecase Service
	Config  *config.Config
}

// New - создание нового сервера.
func New(usecase Service, config *config.Config) *Handler {
	return &Handler{
		usecase: usecase,
		Config:  config,
	}
}

// ListUserData - получение всех данных пользователя.
func (h *Handler) ListUserData(ctx context.Context, _ *emptypb.Empty) (*pb.UserDataResponse, error) {
	var response pb.UserDataResponse

	var (
		err   error
		creds []model.Credentials
		cards []model.PaymentCard
		bins  []model.Binary
		texts []model.Text
	)

	userID := model.UserID("2345423rfsdvf")

	creds, cards, bins, texts, err = h.usecase.GetAll(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	credsPb := convertCreds(creds)
	cardsPb := convertCards(cards)
	binsPb := convertBins(bins)
	textsPb := convertTexts(texts)

	response.SetCredentials(credsPb)
	response.SetPaymentCard(cardsPb)
	response.SetBinary(binsPb)
	response.SetText(textsPb)

	return &response, nil
}

// Credentials - создание данных типа "логин/пароль".
func (h *Handler) AddCredentials(ctx context.Context, in *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	var response pb.AddCredentialsResponse

	userID := model.UserID("2345423rfsdvf")
	// TODO: userID := ctx.Value(constants.UserIDKey).(string)

	credPb := in.GetCredentials()

	cred := model.Credentials{
		UUID:     userID,
		Login:    credPb.GetLogin(),
		Password: credPb.GetPassword(),
		Info:     credPb.GetInfo(),
	}

	err := h.usecase.Create(ctx, &cred)
	if err != nil {
		if errors.Is(err, usecase.ErrDataRegistered) {
			return nil, status.Error(codes.AlreadyExists, `credentials already exist`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	response.SetAlias(cred.Alias)

	return &response, nil
}

// GetCredentials - получение информации о данных типа "логин/пароль".
func (h *Handler) GetCredentials(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	var response pb.GetCredentialsResponse

	userID := model.UserID("2345423rfsdvf")
	// TODO: userID := ctx.Value(constants.UserIDKey).(string)

	alias := in.GetAlias()

	cred := model.Credentials{
		UUID:  userID,
		Alias: alias,
	}

	err := h.usecase.Get(ctx, &cred)
	if err != nil {
		if errors.Is(err, usecase.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, `credentials not found`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	credPb := pb.Credentials{}

	credPb.SetLogin(cred.Login)
	credPb.SetPassword(cred.Password)
	credPb.SetInfo(cred.Info)

	response.SetCredentials(&credPb)

	return &response, nil
}

// DeleteCredentials - удаление данных типа "логин/пароль".
// TODO: не работает
func (h *Handler) DeleteCredentials(ctx context.Context, in *pb.DelCredentialsRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty

	aliases := in.GetAlias()
	creds := make([]any, 0, len(aliases))
	for _, alias := range aliases {
		creds = append(creds, model.Credentials{Alias: alias})
	}

	err := h.usecase.Delete(ctx, creds)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &response, nil
}

// PaymentCard - создание данных типа "банковская карта".
func (h *Handler) PaymentCard(ctx context.Context, in *pb.AddPaymentCardRequest) (*pb.AddPaymentCardResponse, error) {
	var response *pb.AddPaymentCardResponse
	return response, nil
}

// GetPaymentCard - получение данных типа "банковская карта".
func (h *Handler) GetPaymentCard(ctx context.Context, in *pb.GetPaymentCardRequest) (*pb.GetPaymentCardResponse, error) {
	var response *pb.GetPaymentCardResponse
	return response, nil
}

// DeletePaymentCard - удаление данных типа "банковская карта".
func (h *Handler) DeletePaymentCard(ctx context.Context, in *pb.DelPaymentCardRequest) (*emptypb.Empty, error) {
	var response *emptypb.Empty
	return response, nil
}

// Binary - создание данных типа "бинарные данные".
func (h *Handler) Binary(ctx context.Context, in *pb.AddBinaryRequest) (*pb.AddBinaryResponse, error) {
	var response *pb.AddBinaryResponse
	return response, nil
}

// GetBinary - получение данных типа "бинарные данные".
func (h *Handler) GetBinary(ctx context.Context, in *pb.GetBinaryRequest) (*pb.GetBinaryResponse, error) {
	var response *pb.GetBinaryResponse
	return response, nil
}

// DeleteBinary - удаление данных типа "бинарные данные".
func (h *Handler) DeleteBinary(ctx context.Context, in *pb.DelBinaryRequest) (*emptypb.Empty, error) {
	var response *emptypb.Empty
	return response, nil
}

// Text - создание данных типа "текстовые данные".
func (h *Handler) Text(ctx context.Context, in *pb.AddTextRequest) (*pb.AddTextResponse, error) {
	var response *pb.AddTextResponse
	return response, nil
}

// GetText - получение данных типа "текстовые данные".
func (h *Handler) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	var response *pb.GetTextResponse
	return response, nil
}

// DeleteText - удаление данных типа "текстовые данные".
func (h *Handler) DeleteText(ctx context.Context, in *pb.DelTextRequest) (*emptypb.Empty, error) {
	var response *emptypb.Empty
	return response, nil
}
