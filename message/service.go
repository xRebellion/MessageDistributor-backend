package message

import (
	"context"

	"github.com/xRebellion/MessageDistributor/integration"
	"github.com/xRebellion/MessageDistributor/recipient"
)

type Message struct {
	Contents string
	SentTo   []recipient.Contact
}

type Service struct {
	integrations map[string]integration.Integration
}

func (s *Service) SendMessage(ctx context.Context, subject string, message string, destinations []recipient.Contact) []error {
	errors := []error{}
	conversionMap := map[string]*Message{}
	for _, d := range destinations {
		if e, alreadyExists := conversionMap[d.MediaName]; alreadyExists {
			e.SentTo = append(e.SentTo, d)
			continue
		}
		res, err := s.integrations[d.MediaName].ConvertMessage(message)
		if err != nil {
			errors = append(errors, err)
		}
		conversionMap[d.MediaName] = &Message{
			Contents: res,
			SentTo:   []recipient.Contact{d},
		}
	}
	if len(errors) > 0 {
		return errors
	}
	for mediaName, message := range conversionMap {
		mediaUserIDs := []string{}
		for _, dest := range message.SentTo {
			mediaUserIDs = append(mediaUserIDs, dest.MediaUserID)
		}
		errs := s.integrations[mediaName].SendMessage(ctx, subject, message.Contents, mediaUserIDs)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}
	if len(errors) > 0 {
		return errors
	}

	return nil
}

func NewService(integrations map[string]integration.Integration) *Service {
	return &Service{
		integrations: integrations,
	}
}
