package grpc

import (
	"github.com/Di-nis/gophKeeper/internal/model"
	pb "github.com/Di-nis/gophKeeper/pkg/proto"
)

// convertCreds - конвертация данных типа "логин/пароль".
func convertCreds(credsPb []*pb.Credentials) []*model.Credentials {
	creds := make([]*model.Credentials, len(credsPb))

	for i, credPb := range credsPb {
		creds[i] = &model.Credentials{
			Login:    credPb.GetLogin(),
			Password: credPb.GetPassword(),
			Info:     credPb.GetInfo(),
		}
	}
	return creds
}

// convertCards - конвертация данных типа "банковская карта".
func convertCards(cardsPb []*pb.PaymentCard) []*model.PaymentCard {
	cards := make([]*model.PaymentCard, len(cardsPb))

	for i, cardPb := range cardsPb {
		cards[i] = &model.PaymentCard{
			Number:   cardPb.GetNumber(),
			ExpMonth: cardPb.GetExpMonth(),
			ExpYear:  cardPb.GetExpYear(),
			CVV:      cardPb.GetCvv(),
			Holder:   cardPb.GetHolder(),
			Info:     cardPb.GetInfo(),
		}
	}
	return cards
}

// convertBinaries - конвертация данных типа "бинарные данные".
func convertBinary(binaryPb []*pb.Binary) []*model.Binary {
	bins := make([]*model.Binary, len(binaryPb))

	for i, binPb := range binaryPb {
		bins[i] = &model.Binary{
			Data: binPb.GetData(),
			Info: binPb.GetInfo(),
		}
	}
	return bins
}

// convertText - конвертация данных типа "текстовые данные".
func convertText(textPb []*pb.Text) []*model.Text {
	texts := make([]*model.Text, len(textPb))

	for i, textPb := range textPb {
		texts[i] = &model.Text{
			Data: textPb.GetData(),
			Info: textPb.GetInfo(),
		}
	}

	return texts
}
