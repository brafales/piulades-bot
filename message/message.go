package message

import (
	"fmt"

	"github.com/brafales/piulades-bot/pinchito"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

//Build returns a list of Chattable objects to send via Telegram
func Build(chatID int64, user string, images [][]byte, text string) ([]tgbotapi.Chattable, error) {
	imageCount := len(images)
	if imageCount > 0 {
		chattables := make([]tgbotapi.Chattable, imageCount+1, imageCount+1)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.DisableWebPagePreview = true
		chattables[0] = msg
		for i, image := range images {
			file := tgbotapi.FileBytes{
				Name:  fmt.Sprintf("Image %d", i),
				Bytes: image,
			}
			msg := tgbotapi.NewPhotoUpload(chatID, file)
			chattables[i+1] = msg
		}
		return chattables, nil
	}
	msg := tgbotapi.NewMessage(chatID, text)
	return []tgbotapi.Chattable{msg}, nil
}

func BuildLog(chatID int64, log pinchito.Log) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(chatID, log.PrettyText())
	return msg
}
