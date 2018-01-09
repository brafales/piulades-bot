package message

import (
	"fmt"

	"github.com/brafales/piulades-bot/pinchito"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
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
	msg := tgbotapi.NewMessage(chatID, log.TelegramText())
	return msg
}

func GetFWFromUsername(message *tgbotapi.Message) string {
	return message.ForwardFrom.UserName
}

func BuildNewLogLine(nick string, message *tgbotapi.Message) string {

	// TODO Properly handle other types of messages
	// (Audio, img, ...)

	line := "[" + time.Unix(int64(message.ForwardDate), 0).Format("15:04:05") + "] "

	if len(nick) > 0 {
		line += "<" + nick + "> "
	} else {
		line += "< ??? > "
	}

	line += message.Text

	if message.Photo != nil {
		line += "[imatge]"
	} else if message.Audio != nil {
		line += "[àudio]"
	} else if message.Document != nil {
		line += "[fitxer]"
	} else if message.Contact != nil {
		line += "[contacte]"
	} else if message.Location != nil {
		line += "[localització]"
	}

	return line + "\n"
}
