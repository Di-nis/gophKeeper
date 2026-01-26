package usecase

import (
	"github.com/Di-nis/gophKeeper/internal/model"
)

type ClientG interface {
	AddCredentials(cred model.Credentials) error
	GetCredentials(alias string) error
	DelCredentials(alias []string) error
}

type ClientH interface {
	// Register(model.Auth) error
	// Login(model.Auth) error
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
func (u *Usecase) Execute() error {
	item := getItem(u.command)

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
