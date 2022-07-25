package implementation

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TelegramIntegration struct {
	Bot                *tgbotapi.BotAPI
	ChatIDCollection   *mongo.Collection
	htmlSanitizePolicy *bluemonday.Policy
}

type TelegramChatID struct {
	PhoneNumber string `bson:"phoneNumber"`
	ChatID      int64  `bson:"chatID"`
}

func TelegramInit(token string, collection *mongo.Collection) *TelegramIntegration {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	p := bluemonday.NewPolicy()
	p.AllowElements("b", "strong", "i", "em", "code", "pre")
	p.AllowAttrs("href").OnElements("a")

	telegramIntegration := &TelegramIntegration{
		Bot:                botAPI,
		ChatIDCollection:   collection,
		htmlSanitizePolicy: p,
	}

	go telegramIntegration.StartRegistrationHandler()
	return telegramIntegration
}

func (ti *TelegramIntegration) StartRegistrationHandler() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := ti.Bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if strings.HasPrefix(update.Message.Command(), "register") {
			phoneNumber := strings.Replace(update.Message.Command(), "register_", "", 1)
			ctx := context.Background() // background context is used because seems like no need of context injection in these parts
			var res TelegramChatID
			err := ti.ChatIDCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&res)
			if err != mongo.ErrNoDocuments {
				log.Printf("[telegram] %s this number is already registered", phoneNumber)
				continue
			}
			_, err = ti.ChatIDCollection.InsertOne(ctx, TelegramChatID{
				PhoneNumber: phoneNumber,
				ChatID:      update.Message.Chat.ID,
			})
			if err != nil {
				log.Fatalf("[telegram] error while registering %s: %v\n", phoneNumber, err)
				continue
			}
			msg.Text = phoneNumber + " has been successfully registered!"
		}

		if _, err := ti.Bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (ti *TelegramIntegration) SendMessage(ctx context.Context, subject string, message string, mediaUserIDs []string) []error {
	message = "<b>" + subject + "</b>\n\n" + message
	errs := []error{}

	for _, phoneNumber := range mediaUserIDs {
		res := TelegramChatID{}
		err := ti.ChatIDCollection.FindOne(ctx, bson.M{"phoneNumber": phoneNumber}).Decode(&res)
		if err != nil {
			errs = append(errs, fmt.Errorf("[telegram] error while sending to %s: %v", phoneNumber, err))
			continue
		}

		msg := tgbotapi.NewMessage(res.ChatID, message)
		msg.ParseMode = tgbotapi.ModeHTML
		if _, err := ti.Bot.Send(msg); err != nil {
			errs = append(errs, fmt.Errorf("[telegram] error while sending to %s: %v", phoneNumber, err))
			continue
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (ti *TelegramIntegration) ConvertMessage(message string) (string, error) {
	return ti.htmlSanitizePolicy.Sanitize(message), nil
}
