package usecase

import (
	"github.com/Di-nis/gophKeeper/internal/model"
)

// getItem - функция для преобразования команды в объект.
func getItem(command model.Command) any {
	switch command.Item {
	case "credentials":
		return model.Credentials{
			Login:    command.Value[0],
			Password: command.Value[1],
			Info:     command.Value[2],
		}
	case "paymendcard":
		return model.PaymentCard{
			Number: command.Value[0],
			// ExpMonth: int32(command.Value[1]),
			// ExpYear: command.Value[2],
			CVV:    command.Value[3],
			Holder: command.Value[4],
			Info:   command.Value[5],
		}
	case "auth":
		return model.Auth{
			Login:    command.Value[0],
			Password: command.Value[1],
		}
	}

	return nil
}
