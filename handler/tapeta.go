package handler

import (
	"regexp"

	"github.com/brafales/piulades-bot/pinchito"

	"github.com/brafales/piulades-bot/message"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Tapeta struct {
	next   Handler
	chatID int64
	bot    *tgbotapi.BotAPI
}

func (t *Tapeta) Handle(update tgbotapi.Update) error {
	match, err := t.match(update.Message.Text)
	if err != nil {
		return err
	}

	if match {
		log := pinchito.Log{}
		telegramMessage := message.BuildLog(t.chatID, log)
		t.bot.Send(telegramMessage)
	}
	return t.next.Handle(update)
}

func (t *Tapeta) match(s string) (bool, error) {
	return regexp.MatchString("^tapeta$", s)
}
