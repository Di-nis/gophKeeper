package grpc

import (
	pb "github.com/Di-nis/gophKeeper/pkg/proto"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// convertCreds - конвертация данных типа "логин/пароль".
func convertCreds(creds []model.Credentials) []*pb.Credentials {
	credsPb := make([]*pb.Credentials, 0, len(creds))

	for _, cred := range creds {
		credPb := pb.Credentials{}
		credPb.SetLogin(cred.Login)
		credPb.SetPassword(cred.Password)
		credPb.SetAlias(cred.Alias)
		credPb.SetInfo(cred.Info)

		credsPb = append(credsPb, &credPb)
	}
	return credsPb
}

// convertCards - конвертация данных типа "банковская карта".
func convertCards(cards []model.PaymentCard) []*pb.PaymentCard {
	cardsPb := make([]*pb.PaymentCard, 0, len(cards))

	for _, card := range cards {
		cardPb := pb.PaymentCard{}
		cardPb.SetNumber(card.Number)
		cardPb.SetExpMonth(card.ExpMonth)
		cardPb.SetExpYear(card.ExpYear)
		cardPb.SetCvv(card.CVV)
		cardPb.SetHolder(card.Holder)
		cardPb.SetAlias(card.Alias)
		cardPb.SetInfo(card.Info)

		cardsPb = append(cardsPb, &cardPb)
	}
	return cardsPb
}

// convertBins - конвертация данных типа "бинарные данные".
func convertBins(bins []model.Binary) []*pb.Binary {
	binsPb := make([]*pb.Binary, 0, len(bins))

	for _, bin := range bins {
		binPb := pb.Binary{}
		binPb.SetData(bin.Data)
		binPb.SetAlias(bin.Alias)
		binPb.SetInfo(bin.Info)

		binsPb = append(binsPb, &binPb)
	}
	return binsPb
}

// convertTexts - конвертация данных типа "текстовые данные".
func convertTexts(texts []model.Text) []*pb.Text {
	textsPb := make([]*pb.Text, 0, len(texts))

	for _, text := range texts {
		textPb := pb.Text{}
		textPb.SetData(text.Data)
		textPb.SetAlias(text.Alias)
		textPb.SetInfo(text.Info)

		textsPb = append(textsPb, &textPb)
	}
	return textsPb
}
