package main

import (
	"github.com/brafales/piulades-bot/configuration"
	"github.com/brafales/piulades-bot/handler"
	"github.com/brafales/piulades-bot/pinchito"
	"github.com/brafales/piulades-bot/twitter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
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

	pinchitoClient := pinchito.Client{
		PinchitoHost: config.PinchitoHost,
	}

	twitterHandler := &handler.Twitter{Bot: bot, ChatID: config.ChatID, TwitterClient: twitterClient}
	tapetaHandler := &handler.Tapeta{Bot: bot, ChatID: config.ChatID, PinchitoClient: pinchitoClient}
	searchHandler := &handler.Search{Bot: bot, ChatID: config.ChatID, PinchitoClient: pinchitoClient}

	crearHandler := &handler.Crear{
		Bot:            bot,
		ChatID:         config.ChatID,
		ActiveLogs:     map[int]*pinchito.PlogData{},
		AuthToken:      config.PinchitoAuthToken,
		PinchitoClient: pinchitoClient,
	}
	handlers := []handler.Handler{twitterHandler, tapetaHandler, searchHandler, crearHandler}

	log.Println("Ready to handle messages")
	for update := range updates {
		for _, h := range handlers {
			log.Println("handling")
			h.Handle(update)
		}
	}
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
