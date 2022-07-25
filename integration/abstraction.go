package integration

import "context"

type Integration interface {
	SendMessage(ctx context.Context, subject string, message string, mediaUserIDs []string) []error
	ConvertMessage(message string) (string, error) // could be merged with SendMessage at the current state, but will leave this for now.
}
