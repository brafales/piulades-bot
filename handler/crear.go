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

const titlePending string = "#Pending Title#"


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


	if t.userHasLogWithPendingTitle(update.Message) {
		return t.handleTitle(update.Message)
	}


	if t.userHasLogInProgress(update.Message) {
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
	return regexp.MatchString("^" + cmdStartLog + "", message.Command())
}

func (t *Crear) matchEndLog(message *tgbotapi.Message) (bool, error) {
	if message == nil {
		return false, nil
	}
	return regexp.MatchString("^" + cmdEndLog + "", message.Command())
}

func (t *Crear) matchInlineKeyboardReply(update tgbotapi.Update) (bool, error) {
	return update.CallbackQuery != nil, nil
}

func (t *Crear) startLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Starting a new Log for UserID ", message.From.ID)

	if t.userHasLogInProgress(message) {
		delete(t.ActiveLogs, message.From.ID)

		tgMessage := tgbotapi.NewMessage(message.Chat.ID, "You already had a log in progress. Too bad. It has been discarded")
		t.Bot.Send(tgMessage)
	}

	t.ActiveLogs[message.From.ID] = &pinchito.Log{}

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Ready to start creating a new Log.\nForward the messages you want to add. Once finished, type /" + cmdEndLog)
	t.Bot.Send(tgMessage)

	return nil
}

func (t *Crear) endLogFromMessage(message *tgbotapi.Message) error {
	//TODO Handle non-started Log

	log.Println("Ending existing log for UserID ", message.From.ID)

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		tgMessage := tgbotapi.NewMessage(message.Chat.ID, "I haven't found any active log. Ignoring you")
		t.Bot.Send(tgMessage)
		return nil
	}

	pinLog.Nota = 0
	//TODO Fill ID, Autor
	logTime := time.Now()
	pinLog.Dia = logTime.Format("02/01/2006")
	pinLog.Hora = logTime.Format("15:04")

	//TODO Ask for Protagonista

	//You can provide the title as an optional argument in cmdEndLog
	titol := message.CommandArguments()
	if len(titol) > 0 {
		pinLog.Titol = titol
		t.sendLogSummary(message)
	} else {
		pinLog.Titol = titlePending
		tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Which title do you want the log to have?")
		t.Bot.Send(tgMessage)
	}


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

func (t *Crear) userHasLogInProgress(message *tgbotapi.Message) bool {
	return t.ActiveLogs[message.From.ID] != nil
}
func (t *Crear) userHasLogWithPendingTitle(message *tgbotapi.Message) bool {
	if message == nil {
		return false
	}

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return false
	}

	return pinLog.Titol == titlePending
}
func (t *Crear) handleTitle(message *tgbotapi.Message) error {
	pinLog := t.ActiveLogs[message.From.ID]
	pinLog.Titol = message.Text

	return t.sendLogSummary(message)
}

func (t *Crear) sendLogSummary(message *tgbotapi.Message) error {
	pinLog := t.ActiveLogs[message.From.ID]

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "I've created the following Log (" + pinLog.Titol + "):")
	t.Bot.Send(tgMessage)
	tgMessage = tgbotapi.NewMessage(message.Chat.ID, pinLog.Text)

	inlineButtons := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Save Log", keybReplySaveLog), tgbotapi.NewInlineKeyboardButtonData("Discard Log", keybReplyDiscardLog))
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineButtons)
	t.Bot.Send(tgMessage)

	return nil
}