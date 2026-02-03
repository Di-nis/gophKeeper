package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
	in "github.com/Di-nis/gophKeeper/internal/client/cli/input"
	out "github.com/Di-nis/gophKeeper/internal/client/cli/output"
	"github.com/Di-nis/gophKeeper/internal/client/config"
	repo "github.com/Di-nis/gophKeeper/internal/client/repository/local"
	"github.com/Di-nis/gophKeeper/internal/model"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
)

var (
	// ErrMethodNotAllowed - ошибка, возникающая при попытке выполнить неизвестный метод.
	ErrMethodNotAllowed = errors.New("error method not allowed")
	// ErrUnknownType - ошибка, возникающая при попытке обработать неизвестный тип.
	ErrUnknownType = errors.New("unknown type")
	// ErrConvertion - ошибка, возникающая при попытке обработать неизвестный тип.
	ErrConvertion = errors.New("error conversion")
)

// Repository - интерфейс для вставки данных в базу данных.
type Repository interface {
	Write(model.Common) error
	Close() error
}

// Client - интерфейс для работы с клиентом.
type ClientG interface {
	AddCredentials(context.Context, model.Credentials) (string, error)
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

	Sync() (model.Common, error)
}

// ClientH - интерфейс для работы с клиентом.
type ClientH interface {
	Register(context.Context, model.Auth) error
	Login(context.Context, model.Auth) error
}

// Usecase - структура Usecase.
type Usecase struct {
	clientGRPC ClientG
	clientHTTP ClientH
	input      in.Input
	Output     out.Output
	repo       Repository
}

// New - создание структуры Usecase.
func New(cfg *config.Config) (*Usecase, error) {
	input := in.New()
	output := out.New()

	httpClient := hc.New(cfg)
	gRPCClient, err := gc.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("gRPC client initialization error:", err)
	}

	repo, err := repo.New(cfg.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("repo initialization error:", err)
	}

	return &Usecase{
		clientGRPC: gRPCClient,
		clientHTTP: httpClient,
		input:      input,
		Output:     output,
		repo:       repo,
	}, nil
}

// Execute - метод для выполнения команд.
func (u *Usecase) Execute(ctx context.Context) error {
	u.input.Parser()

	if err := u.router(ctx); err != nil {
		u.Output.PrintError()
		return err
	}

	u.Output.PrintSuccess()
	return nil
}

// Execute - метод для выполнения команд.
func (u *Usecase) router(ctx context.Context) error {
	valuesOut := make([]string, 0)

	switch u.input.Method {
	case cli.MethodRegister:
		err := u.register(ctx)
		if err != nil {
			return err
		}
	case cli.MethodLogin:
		err := u.login(ctx)
		if err != nil {
			return err
		}
	case cli.MethodAdd:
		alias, err := u.Adder(ctx)
		if err != nil {
			return err
		}
		valuesOut = append(valuesOut, alias)

	case cli.MethodGet:
		values, err := u.Getter()
		if err != nil {
			return errors.ErrUnsupported
		}
		valuesOut = append(valuesOut, values...)
	case cli.MethodDelete:
		err := u.Deleter()
		if err != nil {
			return err
		}
	case cli.MethodSync:
		err := u.sync(ctx)
		if err != nil {
			return err
		}
	default:
		return ErrMethodNotAllowed
	}

	u.Output.SetMethod(u.input.Method).SetItem(u.input.Item).SetValues(valuesOut...)

	return nil
}

// register - регистрация пользователя.
func (u *Usecase) register(ctx context.Context) error {
	auth := model.Auth{}
	auth.SetLogin(u.input.Values[0]).SetPassword(u.input.Values[1])
	return u.clientHTTP.Register(ctx, auth)
}

// login - авторизация пользователя.
func (u *Usecase) login(ctx context.Context) error {
	auth := model.Auth{}
	auth.SetLogin(u.input.Values[0]).SetPassword(u.input.Values[1])
	return u.clientHTTP.Login(ctx, auth)
}

// Adder - общий метод по добавлению данных.
func (u *Usecase) Adder(ctx context.Context) (string, error) {
	switch u.input.Item {
	case cli.ItemCredentials:
		cred, err := getCredentials(u.input.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddCredentials(ctx, cred)
	case cli.ItemPaymentCard:
		card, err := getPaymentCard(u.input.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddPaymentCard(card)
	case cli.ItemBinary:
		bin, err := getBinary(u.input.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddBinary(bin)
	case cli.ItemText:
		text, err := getText(u.input.Values)
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
	alias := u.input.Values[0]

	switch u.input.Item {
	case cli.ItemCredentials:
		var item model.Credentials
		item, err := u.clientGRPC.GetCredentials(alias)
		if err != nil {
			return nil, err
		}

		values = append(values, item.Login, item.Password, item.Info)
	case cli.ItemPaymentCard:
		item, err := u.clientGRPC.GetPaymentCard(alias)
		if err != nil {
			return nil, err
		}

		values = append(values, item.Number, item.ExpMonth, item.ExpYear, item.CVV, item.Holder, item.Info)
	case cli.ItemBinary:
		item, err := u.clientGRPC.GetBinary(alias)
		if err != nil {
			return nil, err
		}

		values = append(values, string(item.Data), item.Info)
	case cli.ItemText:
		item, err := u.clientGRPC.GetText(alias)
		if err != nil {
			return nil, err
		}

		values = append(values, item.Data, item.Info)
	default:
		return nil, ErrUnknownType
	}

	return values, nil
}

// Deleter - общий метод по получению данных.
func (u *Usecase) Deleter() error {
	alias := u.input.Values[0]

	switch u.input.Item {
	case cli.ItemCredentials:
		return u.clientGRPC.DelCredentials(alias)
	case cli.ItemPaymentCard:
		return u.clientGRPC.DelPaymentCard(alias)
	case cli.ItemBinary:
		return u.clientGRPC.DelBinary(alias)
	case cli.ItemText:
		return u.clientGRPC.DelText(alias)
	default:
		return ErrUnknownType
	}

}

// sync - синхронизация данных.
func (u *Usecase) sync(ctx context.Context) error {
	common, err := u.clientGRPC.Sync()
	if err != nil {
		return err
	}

	return u.repo.Write(common)
}
