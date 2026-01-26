package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/Di-nis/gophKeeper/internal/model"
)

const minLength = 4

var ErrLenArgs = errors.New("command line argument length error")

type grpcClient interface {
	AddCredentials(model.Credentials) error 
}


type Cli struct{
	Args []string
	grpcClient 
}

func New(grpcClient) (*Cli, error)  {
	a := os.Args
	fmt.Println(a)
	if len(os.Args) < minLength {
		return nil, ErrLenArgs
	}
	return &Cli{
		Args: os.Args[1:],
	}, nil
}

func (cli *Cli) Router() {
	switch {
	case os.Args[1] == "add" && os.Args[2] == "credentials":
		cred := model.Credentials{
			Login: os.Args[3],
			Password: os.Args[4],
			Info: os.Args[5],
		}
		err := cli.Credentials(cred)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("данные сохранены")
	// case os.Args[1] == "list":
	// 	list()
	// case os.Args[1] == "remove":
	// 	remove()
	}
}
