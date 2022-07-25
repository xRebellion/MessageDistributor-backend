package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xRebellion/MessageDistributor/integration"
	integrationImpl "github.com/xRebellion/MessageDistributor/integration/implementation"
	"github.com/xRebellion/MessageDistributor/message"
	"github.com/xRebellion/MessageDistributor/recipient"
	"github.com/xRebellion/MessageDistributor/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	bgctx := context.Background()
	dbClient, err := mongo.Connect(bgctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := dbClient.Database("MessageDistributor")

	discordHttpClient := http.Client{}
	integrationRegistry := map[string]integration.Integration{
		"email":    integrationImpl.EmailInit(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET")),
		"discord":  integrationImpl.DiscordInit("http://localhost:7070", discordHttpClient),
		"telegram": integrationImpl.TelegramInit(os.Getenv("TELEGRAM_SECRET"), db.Collection("telegramChatID")),
	}
	recipientService := recipient.NewService(integrationRegistry, db.Collection("recipients"))
	recipientController := recipient.NewController(recipientService)
	messageService := message.NewService(integrationRegistry)
	messageController := message.NewController(messageService)
	router := router.NewChiRouter(
		recipientController,
		messageController,
	).InitRouter()
	fmt.Printf("Running at localhost:8080\n\n")
	log.Fatal(http.ListenAndServe("localhost:8080", router))

}
