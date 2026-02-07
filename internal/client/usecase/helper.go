package usecase

import (
	// in "github.com/Di-nis/gophKeeper/internal/client/cli/input"
	"github.com/Di-nis/gophKeeper/internal/model"
)

// //
// func getItemAdd(c in.Input) (any, error) {
// 	switch c.Item {
// 	case "credentials":
// 		cred := model.Credentials{
// 			Login:    c.Values[0],
// 			Password: c.Values[1],
// 		}
// 		return cred, nil
// 	default:
// 		return nil, ErrMethodNotAllowed
// 	}
// }

// getCredentials - формирование экземпляра Credentials.
func getCredentials(values []string) (model.Credentials, error) {
	if len(values) != 3 {
		return model.Credentials{}, ErrConvertion
	}
	return model.Credentials{
		Login:    values[0],
		Password: values[1],
		Info:     values[2],
	}, nil
}

// getPaymentCard - формирование экземпляра PaymentCard.
func getPaymentCard(values []string) (model.PaymentCard, error) {
	if len(values) != 6 {
		return model.PaymentCard{}, ErrConvertion
	}

	return model.PaymentCard{
		Number:   values[0],
		ExpMonth: values[1],
		ExpYear:  values[2],
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
