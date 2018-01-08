package configuration

import (
	"errors"
	"os"
	"strconv"
)

//Configuration holds application configuration
type Configuration struct {
	TwitterAPIKey    string
	TwitterAPISecret string
	ChatID           int64
	Port             string
	BotKey           string
	CallbackURL      string
}

//New returns a new Config
func New() (*Configuration, error) {
	chatID, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		return &Configuration{}, errors.New("Could not get chat id from environment")
	}
	return &Configuration{
		TwitterAPIKey:    os.Getenv("TWITTER_API_KEY"),
		TwitterAPISecret: os.Getenv("TWITTER_API_SECRET"),
		Port:             os.Getenv("PORT"),
		BotKey:           os.Getenv("BOT_KEY"),
		CallbackURL:      os.Getenv("CALLBACK_URL"),
		ChatID:           chatID,
	}, nil
}
