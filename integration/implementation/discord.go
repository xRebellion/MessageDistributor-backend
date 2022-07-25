package implementation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type DiscordIntegration struct {
	apiBaseURL string
	httpClient http.Client
}

func DiscordInit(baseURL string, client http.Client) *DiscordIntegration {
	return &DiscordIntegration{
		apiBaseURL: baseURL,
		httpClient: client,
	}
}

type DiscordRequest struct {
	UserID  string `json:"discordUserId"`
	Message string `json:"message"`
}

// SendMessage -- message is in Markdown
func (di *DiscordIntegration) SendMessage(ctx context.Context, subject string, message string, mediaUserIDs []string) []error {
	message = "**" + subject + "**\n\n" + message
	errs := []error{}
	for _, userID := range mediaUserIDs {
		payload := DiscordRequest{
			UserID:  userID,
			Message: message,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			errs = append(errs, fmt.Errorf("[discord] error while sending to ID: %s, %v", userID, err))
			continue
		}
		request, err := http.NewRequest("POST", di.apiBaseURL+"/send", bytes.NewBuffer(jsonPayload))
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			errs = append(errs, fmt.Errorf("[discord] error while sending to ID: %s, %v", userID, err))
			continue
		}
		response, err := di.httpClient.Do(request)
		if err != nil {
			errs = append(errs, fmt.Errorf("[discord] error while sending to ID: %s, %v", userID, err))
			continue
		}
		fmt.Print(response.Body)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ConvertMessage converts message from html to markdown
func (di *DiscordIntegration) ConvertMessage(message string) (string, error) {
	markdownConverter := md.NewConverter("", true, nil)
	markdownConverter.AddRules(md.Rule{
		Filter: []string{"u"},
		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
			return md.String("__" + content + "__")
		},
	})
	markdownMessage, err := markdownConverter.ConvertString(message)
	if err != nil {
		return "", err
	}

	return markdownMessage, nil
}
