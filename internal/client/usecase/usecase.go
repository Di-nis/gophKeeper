package usecase

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
	in "github.com/Di-nis/gophKeeper/internal/client/cli/input"
	out "github.com/Di-nis/gophKeeper/internal/client/cli/output"
	"github.com/Di-nis/gophKeeper/internal/client/config"
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

// Repository - интерфейс для вставки данных в базу данных.
type Repository interface {
	Write(model.Common) error
	Close() error
}

// Client - интерфейс для работы с клиентом.
type ClientG interface {
	AddCredentials(context.Context, model.Credentials) (string, error)
	GetCredentials(context.Context, string) (model.Credentials, error)
	DelCredentials(context.Context, string) error

	AddPaymentCard(context.Context, model.PaymentCard) (string, error)
	GetPaymentCard(context.Context, string) (model.PaymentCard, error)
	DelPaymentCard(context.Context, string) error

	AddBinary(context.Context, model.Binary) (string, error)
	GetBinary(context.Context, string) (model.Binary, error)
	DelBinary(context.Context, string) error

	AddText(context.Context, model.Text) (string, error)
	GetText(context.Context, string) (model.Text, error)
	DelText(context.Context, string) error

	Sync(context.Context) (model.Common, error)
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
func New(cfg *config.Config, httpClient ClientH, grpcClient ClientG, repo Repository) (*Usecase, error) {
	input := in.New()
	output := out.New()

	return &Usecase{
		clientHTTP: httpClient,
		clientGRPC: grpcClient,
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
		values, err := u.Getter(ctx)
		if err != nil {
			return errors.ErrUnsupported
		}
		valuesOut = append(valuesOut, values...)
	case cli.MethodDelete:
		err := u.Deleter(ctx)
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
		return u.clientGRPC.AddPaymentCard(ctx, card)
	case cli.ItemBinary:
		bin, err := getBinary(u.input.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddBinary(ctx, bin)
	case cli.ItemText:
		text, err := getText(u.input.Values)
		if err != nil {
			return "", err
		}
		return u.clientGRPC.AddText(ctx, text)
	default:
		return "", ErrUnknownType
	}
}

// Getter - общий метод по получению данных.
func (u *Usecase) Getter(ctx context.Context) ([]string, error) {
	values := make([]string, 0)
	alias := u.input.Values[0]

	switch u.input.Item {
	case cli.ItemCredentials:
		var item model.Credentials
		item, err := u.clientGRPC.GetCredentials(ctx, alias)
		if err != nil {
			return nil, err
		}

		values = append(values, item.Login, item.Password, item.Info)
	case cli.ItemPaymentCard:
		item, err := u.clientGRPC.GetPaymentCard(ctx, alias)
		if err != nil {
			return nil, err
		}

		values = append(values, item.Number, item.ExpMonth, item.ExpYear, item.CVV, item.Holder, item.Info)
	case cli.ItemBinary:
		item, err := u.clientGRPC.GetBinary(ctx, alias)
		if err != nil {
			return nil, err
		}

		values = append(values, string(item.Data), item.Info)
	case cli.ItemText:
		item, err := u.clientGRPC.GetText(ctx, alias)
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
func (u *Usecase) Deleter(ctx context.Context) error {
	alias := u.input.Values[0]

	switch u.input.Item {
	case cli.ItemCredentials:
		return u.clientGRPC.DelCredentials(ctx, alias)
	case cli.ItemPaymentCard:
		return u.clientGRPC.DelPaymentCard(ctx, alias)
	case cli.ItemBinary:
		return u.clientGRPC.DelBinary(ctx, alias)
	case cli.ItemText:
		return u.clientGRPC.DelText(ctx, alias)
	default:
		return ErrUnknownType
	}

}

// sync - синхронизация данных.
func (u *Usecase) sync(ctx context.Context) error {
	common, err := u.clientGRPC.Sync(ctx)
	if err != nil {
		return err
	}

	return u.repo.Write(common)
}
