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

func BuildNewLogLine(message *tgbotapi.Message) string {
	var author string
	if message.ForwardFrom != nil {
		user, err := pinchito.GetUserFromTelegramUsername(message.ForwardFrom.UserName)
		author = "[" + time.Unix(int64(message.ForwardDate), 0).Format("15:04:05") + "] "
		if err == nil {
			author += "<" + user.PinNick + "> "
		}else if len(message.ForwardFrom.UserName) > 0 {
			author += "<" + message.ForwardFrom.UserName + "> "
		} else {
			author += "<" + message.ForwardFrom.FirstName + " " + message.ForwardFrom.LastName + "> "
		}
	} else {
		author = "[" + time.Unix(int64(message.Date), 0).Format("15:04:05") + "] < ??? > "
	}

	return author + message.Text + "\n"
}
