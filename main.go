package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brafales/piulades-bot/configuration"
	"github.com/brafales/piulades-bot/message"
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

	for update := range updates {
		err = processUpdate(update, bot, twitterClient, config.ChatID)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func processUpdate(update tgbotapi.Update,
	bot *tgbotapi.BotAPI,
	twitterClient *twitter.Client,
	chatID int64) error {
	statusID, err := twitter.GetStatusID(update.Message.Text)
	if err != nil {
		return err
	}

	tweet, err := twitterClient.GetTwit(statusID)
	if err != nil {
		return err
	}

	images, err := tweet.ExtendedEntities()
	if err != nil {
		return err
	}
	messages, err := message.Build(chatID,
		update.Message.From.UserName,
		images,
		tweet.PrintableText(update.Message.From.UserName))
	if err != nil {
		return err
	}

	for _, message := range messages {
		_, err := bot.Send(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
