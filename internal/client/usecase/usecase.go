package usecase

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
	"github.com/Di-nis/gophKeeper/internal/model"
)

var (
	// ErrMethodNotAllowed - ошибка, возникающая при попытке выполнить неизвестный метод.
	ErrMethodNotAllowed = errors.New("error method not allowed")
	// ErrUnknownType - ошибка, возникающая при попытке обработать неизвестный тип.
	ErrUnknownType = errors.New("unknown type")
	// ErrConvertion - ошибка, возникающая при попытке обработать неизвестный тип.
	ErrConvertion = errors.New("error conversion")
)

// Client - интерфейс для работы с клиентом.
type ClientG interface {
	AddCredentials(model.Credentials) (string, error)
	GetCredentials(string) (model.Credentials, error)
	DelCredentials(string) error

	AddPaymentCard(model.PaymentCard) (string, error)
	GetPaymentCard(string) (model.PaymentCard, error)
	DelPaymentCard(string) error

	AddBinary(model.Binary) (string, error)
	GetBinary(string) (model.Binary, error)
	DelBinary(string) error

	AddText(model.Text) (string, error)
	GetText(string) (model.Text, error)
	DelText(string) error
}

// ClientH - интерфейс для работы с клиентом.
type ClientH interface {
	Register(context.Context, model.Auth) error
	// Login(context.Context, model.Auth) error
}

// Usecase - структура Usecase.
type Usecase struct {
	clientGRPC ClientG
	clientHTTP ClientH
	cli        cli.Cli
}

// New - создание структуры Usecase.
func New(clientG ClientG, clientH ClientH, c cli.Cli) *Usecase {
	return &Usecase{
		clientGRPC: clientG,
		clientHTTP: clientH,
		cli:        c,
	}
}

// Execute - метод для выполнения команд.
func (u *Usecase) Execute() error {
	switch u.cli.Method {
	case cli.MethodAdd:
		alias, err := u.Adder()
		if err != nil {
			return err
		}
		u.cli.Alias = alias

	case cli.MethodGet:
		values, err := u.Getter()
		if err != nil {
			return errors.ErrUnsupported
		}
		u.cli.Values = values
	// case cli.MethodDel:
	// 	return u.Deleter()
	default:
		return ErrMethodNotAllowed
	}

	u.cli.Print()
	return nil
}

// Adder - общий метод по добавлению данных.
func (u *Usecase) Adder() (string, error) {
	switch u.cli.Item {
	case cli.ItemCredentials:
		cred, err := getCredentials(u.cli.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddCredentials(cred)
	case cli.ItemPaymentCard:
		card, err := getPaymentCard(u.cli.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddPaymentCard(card)
	case cli.ItemBinary:
		bin, err := getBinary(u.cli.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddBinary(bin)
	case cli.ItemText:
		text, err := getText(u.cli.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddText(text)
	default:
		return "", ErrUnknownType
	}
}

// Getter - общий метод по получению данных.
func (u *Usecase) Getter() ([]string, error) {
	values := make([]string, 0)

	switch u.cli.Item {
	case cli.ItemCredentials:
		alias := u.cli.Values[0]
		item, err := u.clientGRPC.GetCredentials(alias)
		if err != nil {
			return nil, err
		}
		values = append(values, item.Login, item.Password)
	// case cli.ItemPaymentCard:
	// 	item, err := u.clientGRPC.GetPaymentCard(u.cli.Alias)
	// case cli.ItemPaymentCard:
	// 	return  u.clientGRPC.GetPaymentCard(u.cli.Alias), nil
	// case cli.ItemBinary:
	// 	return u.clientGRPC.GetBinary(u.cli.Alias), nil
	// case cli.ItemText:
	// 	return u.clientGRPC.GetText(u.cli.Alias), nil
	default:
		return nil, ErrUnknownType
	}

	return values, nil
}
