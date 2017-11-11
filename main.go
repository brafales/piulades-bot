package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brafales/piulades-bot/configuration"
	"github.com/brafales/piulades-bot/handler"
	"github.com/brafales/piulades-bot/twitter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := configuration.New()
	fail(err)

	bot, err := tgbotapi.NewBotAPI(config.BotKey)
	fail(err)

	twitterClient := twitter.NewClient(config.TwitterAPIKey, config.TwitterAPISecret)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.CallbackURL))
	fail(err)

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	twitterHandler := &handler.Twitter{Bot: bot, ChatID: config.ChatID, TwitterClient: twitterClient}
	tapetaHandler := &handler.Tapeta{Bot: bot, ChatID: config.ChatID}
	handlers := []handler.Handler{twitterHandler, tapetaHandler}

	for update := range updates {
		for _, h := range handlers {
			h.Handle(update)
		}
	}
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
