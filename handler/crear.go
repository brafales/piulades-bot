package handler

import (
	"log"
	"regexp"
	"time"

	"github.com/brafales/piulades-bot/pinchito"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Crear struct {
	ChatID     int64
	Bot        *tgbotapi.BotAPI
	ActiveLogs map[int]*pinchito.Log
}

const keybReplySaveLog string = "SaveLog"
const keybReplyDiscardLog string = "DiscardLog"

const cmdStartLog string = "log"
const cmdEndLog string = "end_log"


func (t *Crear) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Crear")

	start, err := t.matchStartLog(update.Message)
	if err != nil {
		return err
	}
	if start {
		return t.startLogFromMessage(update.Message)
	}


	end, err := t.matchEndLog(update.Message)
	if err != nil {
		return err
	}
	if end {
		return t.endLogFromMessage(update.Message)
	}


	isReply, err := t.matchInlineKeyboardReply(update)
	if err != nil {
		return err
	}
	if isReply {
		return t.handleKeyboardReply(update)
	}


	logInProgress, err := t.userHasLogInprogress(update.Message)
	if err != nil {
		return err
	}
	if logInProgress {
		t.appendMessageToLog(update.Message)
	} else {
		tgMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "To Start a new log, type /" + cmdStartLog + " and then forward the messages you want to add. Once finished, type /" + cmdEndLog)
		t.Bot.Send(tgMessage)
	}

	return nil
}

func (t *Crear) matchStartLog(message *tgbotapi.Message) (bool, error) {
	if message == nil {
		return false, nil
	}
	return regexp.MatchString("^/" + cmdStartLog + "", message.Text)
}

func (t *Crear) startLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Starting a new Log for UserID ", message.From.ID)
	//TODO Manage existing log situation

	t.ActiveLogs[message.From.ID] = &pinchito.Log{}

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Ready to start creating a new Log.\nPlease forward the messages you want to add. Once finished, type /" + cmdEndLog)
	t.Bot.Send(tgMessage)

	return nil
}

func (t *Crear) matchEndLog(message *tgbotapi.Message) (bool, error) {
	if message == nil {
		return false, nil
	}
	return regexp.MatchString("^/" + cmdEndLog + "", message.Text)
}

func (t *Crear) endLogFromMessage(message *tgbotapi.Message) error {
	//TODO Handle non-started Log

	log.Println("Ending existing log for UserID ", message.From.ID)

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return nil
	}

	pinLog.Nota = 0
	//TODO Fill ID, Autor, Dia, Hora
	//TODO Ask for Titol, Protagonista

	//log.Println("I've created the following Log:")
	//log.Println(pinLog.Text)

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "I've created the following Log:")
	t.Bot.Send(tgMessage)
	tgMessage = tgbotapi.NewMessage(message.Chat.ID, pinLog.Text)
	//t.Bot.Send(tgMessage)

	//tgMessage = tgbotapi.NewMessage(message.Chat.ID, "Would you like me to save it?")

	inlineButtons := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Save Log", keybReplySaveLog), tgbotapi.NewInlineKeyboardButtonData("Discard Log", keybReplyDiscardLog))
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineButtons)
	t.Bot.Send(tgMessage)

	return nil
}

func (t *Crear) appendMessageToLog(message *tgbotapi.Message) error {
	if message == nil {
		return nil
	}

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return nil
	}

	var author string
	if message.ForwardFrom != nil {
		author = "[" + time.Unix(int64(message.ForwardDate), 0).Format("15:04:05") + "] "
		if len(message.ForwardFrom.UserName) > 0 {
			author += "<" + message.ForwardFrom.UserName + "> "
		} else {
			author += "<" + message.ForwardFrom.FirstName + " " + message.ForwardFrom.LastName + "> "
		}
	} else {
		author = "[" + time.Unix(int64(message.Date), 0).Format("15:04:05") + "] < ??? > "
	}

	pinLog.Text += author + message.Text + "\n"

	return nil
}

func (t *Crear) matchInlineKeyboardReply(update tgbotapi.Update) (bool, error) {
	return update.CallbackQuery != nil, nil
}

func (t *Crear) handleKeyboardReply(update tgbotapi.Update) error {
	if update.CallbackQuery.Data == keybReplySaveLog {
		t.saveLog(update)
	} else if update.CallbackQuery.Data == keybReplyDiscardLog {
		t.discardLog(update)
	}


	return nil
}

func (t *Crear) saveLog(update tgbotapi.Update) error {

	//TODO Save to DB

	delete(t.ActiveLogs, update.CallbackQuery.From.ID)

	oldMessage := update.CallbackQuery.Message
	editMessage := tgbotapi.NewEditMessageText(oldMessage.Chat.ID, oldMessage.MessageID, oldMessage.Text + "\n\nLog saved: http://go.pinchito.com/1427")
	t.Bot.Send(editMessage)

	return nil
}

func (t *Crear) discardLog(update tgbotapi.Update) error {
	delete(t.ActiveLogs, update.CallbackQuery.From.ID)

	oldMessage := update.CallbackQuery.Message
	editMessage := tgbotapi.NewEditMessageText(oldMessage.Chat.ID, oldMessage.MessageID, "Log discarded")
	t.Bot.Send(editMessage)

	return nil
}
func (t *Crear) userHasLogInprogress(message *tgbotapi.Message) (bool, error) {
	return t.ActiveLogs[message.From.ID] != nil, nil
}