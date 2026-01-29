package usecase

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// errMethodNotAllowed - ошибка, возникающая при попытке выполнить неизвестный метод.
var errMethodNotAllowed = errors.New("error method not allowed")

type ClientG interface {
	AddCredentials(cred model.Credentials) error
	GetCredentials(alias string) error
	DelCredentials(alias []string) error
}

type ClientH interface {
	Register(context.Context, model.Auth) error
	Login(context.Context, model.Auth) error
}

type Usecase struct {
	ClientGRPC ClientG
	ClientHTTP ClientH
	command    model.Command
}

func New(clientG ClientG, clientH ClientH, command model.Command) *Usecase {
	return &Usecase{
		ClientGRPC: clientG,
		ClientHTTP: clientH,
		command:    command,
	}
}

// Execute - метод для выполнения команд.
func (u *Usecase) Execute(ctx context.Context) error {
	item := getItem(u.command)

	a, ok := item.(model.Auth)
	if ok {
		switch u.command.Method {
		case "register":
			return u.ClientHTTP.Register(ctx, a)
		case "login":
			return u.ClientHTTP.Login(ctx, a)
		default:
			return errMethodNotAllowed
		}
	}

	v, ok := item.(model.Credentials)
	if ok {
		switch u.command.Method {
		case "add":
			return u.ClientGRPC.AddCredentials(v)
		case "get":
			return u.ClientGRPC.GetCredentials(v.Alias)
		case "del":
			return u.ClientGRPC.DelCredentials([]string{v.Alias})
		}
	}
	return nil
}
