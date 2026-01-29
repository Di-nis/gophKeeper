package usecase

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// ErrMethodNotAllowed - ошибка, возникающая при попытке выполнить неизвестный метод.
var ErrMethodNotAllowed = errors.New("error method not allowed")

// Client - интерфейс для работы с клиентом.
type ClientG interface {
	AddCredentials(cred model.Credentials) error
	GetCredentials(alias string) error
	DelCredentials(alias []string) error
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
	command    model.Command
}

// New - создание структуры Usecase.
func New(clientG ClientG, clientH ClientH, command model.Command) *Usecase {
	return &Usecase{
		clientGRPC: clientG,
		clientHTTP: clientH,
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
			return u.clientHTTP.Register(ctx, a)
		case "login":
			return u.clientHTTP.Login(ctx, a)
		default:
			return ErrMethodNotAllowed
		}
	}

	v, ok := item.(model.Credentials)
	if ok {
		switch u.command.Method {
		case "add":
			return u.clientGRPC.AddCredentials(v)
		case "get":
			return u.clientGRPC.GetCredentials(v.Alias)
		case "del":
			return u.clientGRPC.DelCredentials([]string{v.Alias})
		}
	}
	return nil
}
