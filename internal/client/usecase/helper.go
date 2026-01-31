package usecase

import (
	"strconv"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
	"github.com/Di-nis/gophKeeper/internal/model"
)

func getItemAdd(c cli.Cli) (any, error) {
	switch c.Item {
	case "credentials":
		cred := model.Credentials{
			Login:    c.Values[0],
			Password: c.Values[1],
		}
		return cred, nil
	default:
		return nil, ErrMethodNotAllowed
	}
}

// getCredentials - формирование экземпляра Credentials.
func getCredentials(values []string) (model.Credentials, error) {
	if len(values) != 2 {
		return model.Credentials{}, ErrConvertion
	}
	return model.Credentials{
		Login:    values[0],
		Password: values[1],
	}, nil
}

// getPaymentCard - формирование экземпляра PaymentCard.
func getPaymentCard(values []string) (model.PaymentCard, error) {
	if len(values) != 6 {
		return model.PaymentCard{}, ErrConvertion
	}

	expMonth, err := strconv.ParseInt(values[1], 10, 32)
	if err != nil {
		return model.PaymentCard{}, ErrConvertion
	}
	expYear, err := strconv.ParseInt(values[2], 10, 32)
	if err != nil {
		return model.PaymentCard{}, ErrConvertion
	}

	return model.PaymentCard{
		Number:   values[0],
		ExpMonth: int32(expMonth),
		ExpYear:  int32(expYear),
		CVV:      values[3],
		Holder:   values[4],
		Info:     values[5],
	}, nil
}

// getBinary - формирование экземпляра Binary.
func getBinary(values []string) (model.Binary, error) {
	if len(values) != 2 {
		return model.Binary{}, ErrConvertion
	}
	return model.Binary{
		Data: []byte(values[0]),
		Info: values[1],
	}, nil
}

// getText - формирование экземпляра Text.
func getText(values []string) (model.Text, error) {
	if len(values) != 2 {
		return model.Text{}, ErrConvertion
	}
	return model.Text{
		Data: values[0],
		Info: values[1],
	}, nil
}
