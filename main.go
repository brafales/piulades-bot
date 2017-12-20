package main

import (
	"log"
	"time"
	//"net/http"
	//"os"

	"github.com/brafales/piulades-bot/configuration"
	"github.com/brafales/piulades-bot/handler"
	"github.com/brafales/piulades-bot/pinchito"
	//"github.com/brafales/piulades-bot/twitter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := configuration.New()
	fail(err)

	bot, err := tgbotapi.NewBotAPI(config.BotKey)
	fail(err)
	/*
		twitterClient := twitter.NewClient(config.TwitterAPIKey, config.TwitterAPISecret)

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.CallbackURL))
		fail(err)

		updates := bot.ListenForWebhook("/")
		go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

		twitterHandler := &handler.Twitter{Bot: bot, ChatID: config.ChatID, TwitterClient: twitterClient}
		tapetaHandler := &handler.Tapeta{Bot: bot, ChatID: config.ChatID}
		searchHandler := &handler.Search{Bot: bot, ChatID: config.ChatID}
		handlers := []handler.Handler{twitterHandler, tapetaHandler, searchHandler}

		log.Println("Ready to handle messages")
		for update := range updates {
			for _, h := range handlers {
				log.Println("handling")
				h.Handle(update)
			}
		}
	*/
	crearHandler := &handler.Crear{Bot: bot, ChatID: config.ChatID, ActiveLogs: map[int]*pinchito.Log{}}
	conf := tgbotapi.NewUpdate(0)
	for {
		//log.Println("Conf offset: ", conf.Offset)
		updates, _ := bot.GetUpdates(conf)
		if len(updates) > 0 {
			for _, update := range updates {
				//log.Println("UpdateID: ", update.UpdateID)
				crearHandler.Handle(update)

				if update.UpdateID >= conf.Offset {
					conf.Offset = update.UpdateID + 1
					//log.Println("Conf: ", conf)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
