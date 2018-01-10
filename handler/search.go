package handler

import (
	"log"
	"regexp"

	"github.com/brafales/piulades-bot/pinchito"

	"github.com/brafales/piulades-bot/message"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Search struct {
	ChatID         int64
	Bot            *tgbotapi.BotAPI
	PinchitoClient pinchito.Client
}

func (s *Search) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Search")
	if update.Message == nil {
		return nil
	}
	term := s.match(update.Message.Text)
	if term == "" {
		log.Println("No Search command found")
		return nil
	}

	log, err := s.PinchitoClient.Search(term)
	if err != nil {
		return err
	}

	telegramMessage := message.BuildLog(s.ChatID, log)
	s.Bot.Send(telegramMessage)

	return nil
}

func (search *Search) match(s string) string {
	re := regexp.MustCompile("^/search ([\\w\\s\\.,-\\?!]+)$")
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return ""
	}
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
