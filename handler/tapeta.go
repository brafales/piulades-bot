package handler

import (
	"log"
	"regexp"

	"github.com/brafales/piulades-bot/pinchito"

	"github.com/brafales/piulades-bot/message"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Tapeta struct {
	ChatID int64
	Bot    *tgbotapi.BotAPI
}

func (t *Tapeta) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Tapeta")
	match, err := t.match(update.Message.Text)
	if err != nil {
		return err
	}

	if !match {
		log.Println("No Tapeta command found")
		return nil
	}

	log, err := pinchito.Tapeta()
	if err != nil {
		return err
	}
	telegramMessage := message.BuildLog(t.ChatID, log)
	t.Bot.Send(telegramMessage)

	return nil
}

func (t *Tapeta) match(s string) (bool, error) {
	return regexp.MatchString("^/tapeta$", s)
}
