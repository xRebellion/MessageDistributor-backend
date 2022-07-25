package implementation

import (
	"context"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
)

type EmailIntegration struct {
	emailClient    *mailjet.Client
	defaultSubject string
}

func EmailInit(apiKey string, secretKey string) *EmailIntegration {
	return &EmailIntegration{
		emailClient:    mailjet.NewMailjetClient(apiKey, secretKey),
		defaultSubject: "Mail from MessageDistributor",
	}
}

// SendMessage
func (ei *EmailIntegration) SendMessage(ctx context.Context, subject string, message string, mediaUserIDs []string) []error {
	if subject == "" {
		subject = ei.defaultSubject
	}
	messageInfo := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: "danielyg2904@gmail.com",
			Name:  "Daniel",
		},
		Subject:  subject,
		HTMLPart: message,
	}
	recipients := []mailjet.RecipientV31{}
	for _, email := range mediaUserIDs {
		recipients = append(recipients, mailjet.RecipientV31{Email: email})
	}
	messageInfo.To = (*mailjet.RecipientsV31)(&recipients)

	messages := mailjet.MessagesV31{Info: []mailjet.InfoMessagesV31{messageInfo}}
	_, err := ei.emailClient.SendMailV31(&messages)
	if err != nil {
		return []error{fmt.Errorf("[email] error while sending to address %v: %v", mediaUserIDs, err)}
	}
	return nil
}

// ConvertMessage does not convert anything - uses html directly
func (ei *EmailIntegration) ConvertMessage(message string) (string, error) {
	return message, nil
}
