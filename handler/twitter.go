package handler

import (
	"log"

	"github.com/brafales/piulades-bot/message"
	"github.com/brafales/piulades-bot/twitter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Twitter struct {
	ChatID        int64
	Bot           *tgbotapi.BotAPI
	TwitterClient *twitter.Client
}

func (t *Twitter) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Twitter")
	statusID, err := twitter.GetStatusID(update.Message.Text)
	if err != nil {
		log.Println("No twitter link found")
		return err
	}

	tweet, err := t.TwitterClient.GetTwit(statusID)
	if err != nil {
		return err
	}

	images, err := tweet.ExtendedEntities()
	if err != nil {
		return err
	}
	messages, err := message.Build(t.ChatID,
		update.Message.From.UserName,
		images,
		tweet.PrintableText(update.Message.From.UserName))
	if err != nil {
		return err
	}

	for _, message := range messages {
		_, err := t.Bot.Send(message)
		if err != nil {
			return err
		}
	}
	return nil
}
