// Package grpc - реализация gRPC-сервера.
package grpc

import (
	"context"
	"errors"

	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/auth"
	"github.com/Di-nis/gophKeeper/internal/server/config"
	"github.com/Di-nis/gophKeeper/internal/server/repository"

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
	Delete(context.Context, any) error
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

	userID := ctx.Value(auth.Key).(model.UserID)

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

// AddCredentials - создание данных типа "логин/пароль".
func (h *Handler) AddCredentials(ctx context.Context, in *pb.AddCredentialsRequest) (*pb.AddCredentialsResponse, error) {
	var response pb.AddCredentialsResponse

	userID := ctx.Value(auth.Key).(model.UserID)

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

	userID := ctx.Value(auth.Key).(model.UserID)
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
func (h *Handler) DelCredentials(ctx context.Context, in *pb.DelCredentialsRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty

	userID := ctx.Value(auth.Key).(model.UserID)
	cred := model.Credentials{
		UUID:  userID,
		Alias: in.GetAlias(),
	}

	err := h.usecase.Delete(ctx, cred)
	if err != nil && errors.Is(err, repository.ErrDataNotFound) {
		return nil, status.Error(codes.NotFound, "data not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &response, nil
}

// AddPaymentCard - создание данных типа "банковская карта".
func (h *Handler) AddPaymentCard(ctx context.Context, in *pb.AddPaymentCardRequest) (*pb.AddPaymentCardResponse, error) {
	var response pb.AddPaymentCardResponse

	userID := ctx.Value(auth.Key).(model.UserID)

	cardPb := in.GetCard()

	card := model.PaymentCard{
		UUID:     userID,
		Number:   cardPb.GetNumber(),
		ExpMonth: cardPb.GetExpMonth(),
		ExpYear:  cardPb.GetExpYear(),
		CVV:      cardPb.GetCvv(),
		Info:     cardPb.GetInfo(),
	}

	err := h.usecase.Create(ctx, &card)
	if err != nil {
		if errors.Is(err, usecase.ErrDataRegistered) {
			return nil, status.Error(codes.AlreadyExists, `payment card already exist`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	response.SetAlias(card.Alias)

	return &response, nil
}

// GetPaymentCard - получение данных типа "банковская карта".
func (h *Handler) GetPaymentCard(ctx context.Context, in *pb.GetPaymentCardRequest) (*pb.GetPaymentCardResponse, error) {
	var response pb.GetPaymentCardResponse

	userID := ctx.Value(auth.Key).(model.UserID)
	alias := in.GetAlias()

	card := model.PaymentCard{
		UUID:  userID,
		Alias: alias,
	}

	err := h.usecase.Get(ctx, &card)
	if err != nil {
		if errors.Is(err, usecase.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, `payment card not found`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	cardPb := pb.PaymentCard{}

	cardPb.SetNumber(card.Number)
	cardPb.SetExpMonth(card.ExpMonth)
	cardPb.SetExpYear(card.ExpYear)
	cardPb.SetCvv(card.CVV)
	cardPb.SetHolder(card.Holder)
	cardPb.SetInfo(card.Info)

	response.SetCard(&cardPb)

	return &response, nil
}

// DeletePaymentCard - удаление данных типа "банковская карта".
func (h *Handler) DelPaymentCard(ctx context.Context, in *pb.DelPaymentCardRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty

	userID := ctx.Value(auth.Key).(model.UserID)
	card := model.PaymentCard{
		UUID:  userID,
		Alias: in.GetAlias(),
	}

	err := h.usecase.Delete(ctx, card)
	if err != nil && errors.Is(err, repository.ErrDataNotFound) {
		return nil, status.Error(codes.NotFound, "data not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &response, nil
}

// AddBinary - создание данных типа "бинарные данные".
func (h *Handler) AddBinary(ctx context.Context, in *pb.AddBinaryRequest) (*pb.AddBinaryResponse, error) {
	var response pb.AddBinaryResponse

	userID := ctx.Value(auth.Key).(model.UserID)

	binPb := in.GetBin()

	bin := model.Binary{
		UUID: userID,
		Data: binPb.GetData(),
		Info: binPb.GetInfo(),
	}

	err := h.usecase.Create(ctx, &bin)
	if err != nil {
		if errors.Is(err, usecase.ErrDataRegistered) {
			return nil, status.Error(codes.AlreadyExists, `binary data already exist`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	response.SetAlias(bin.Alias)

	return &response, nil
}

// GetBinary - получение данных типа "бинарные данные".
func (h *Handler) GetBinary(ctx context.Context, in *pb.GetBinaryRequest) (*pb.GetBinaryResponse, error) {
	var response pb.GetBinaryResponse

	userID := ctx.Value(auth.Key).(model.UserID)
	alias := in.GetAlias()

	bin := model.Binary{
		UUID:  userID,
		Alias: alias,
	}

	err := h.usecase.Get(ctx, &bin)
	if err != nil {
		if errors.Is(err, usecase.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, `binary data not found`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	binPb := pb.Binary{}

	binPb.SetData(bin.Data)
	binPb.SetInfo(bin.Info)

	response.SetBin(&binPb)

	return &response, nil

}

// DelBinary - удаление данных типа "бинарные данные".
func (h *Handler) DelBinary(ctx context.Context, in *pb.DelBinaryRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty

	userID := ctx.Value(auth.Key).(model.UserID)
	bin := model.Binary{
		UUID:  userID,
		Alias: in.GetAlias(),
	}

	err := h.usecase.Delete(ctx, bin)
	if err != nil && errors.Is(err, repository.ErrDataNotFound) {
		return nil, status.Error(codes.NotFound, "data not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &response, nil
}

// Text - создание данных типа "текстовые данные".
func (h *Handler) AddText(ctx context.Context, in *pb.AddTextRequest) (*pb.AddTextResponse, error) {
	var response pb.AddTextResponse

	userID := ctx.Value(auth.Key).(model.UserID)

	textPb := in.GetText()

	text := model.Text{
		UUID: userID,
		Data: textPb.GetData(),
		Info: textPb.GetInfo(),
	}

	err := h.usecase.Create(ctx, &text)
	if err != nil {
		if errors.Is(err, usecase.ErrDataRegistered) {
			return nil, status.Error(codes.AlreadyExists, `text data already exist`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	response.SetAlias(text.Alias)

	return &response, nil
}

// GetText - получение данных типа "текстовые данные".
func (h *Handler) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	var response pb.GetTextResponse

	userID := ctx.Value(auth.Key).(model.UserID)
	alias := in.GetAlias()

	text := model.Text{
		UUID:  userID,
		Alias: alias,
	}

	err := h.usecase.Get(ctx, &text)
	if err != nil {
		if errors.Is(err, usecase.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, `text data not found`)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	textPb := pb.Text{}

	textPb.SetData(text.Data)
	textPb.SetInfo(text.Info)

	response.SetText(&textPb)

	return &response, nil
}

// DeleteText - удаление данных типа "текстовые данные".
func (h *Handler) DelText(ctx context.Context, in *pb.DelTextRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty

	userID := ctx.Value(auth.Key).(model.UserID)
	text := model.Text{
		UUID:  userID,
		Alias: in.GetAlias(),
	}

	err := h.usecase.Delete(ctx, text)
	if err != nil && errors.Is(err, repository.ErrDataNotFound) {
		return nil, status.Error(codes.NotFound, "data not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &response, nil
}
