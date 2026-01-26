package usecase

import (
	"github.com/Di-nis/gophKeeper/internal/model")


func getItem(command model.Command) any {
	switch command.Item {
	case "Credentials":
		return model.Credentials{
			Login: command.Value[0],
			Password: command.Value[1],
			Info: command.Value[2],
		}
	case "PaymendCard":
		return model.PaymentCard{
			Number: command.Value[0],
			// ExpMonth: int32(command.Value[1]),
			// ExpYear: command.Value[2],
			CVV: command.Value[3],
			Holder: command.Value[4],
			Info: command.Value[5],
		}
	}
	return nil
}