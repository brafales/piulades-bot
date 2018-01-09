package handler

import (
	"fmt"
	"log"
	"regexp"
	"time"

	pinmessage "github.com/brafales/piulades-bot/message"
	"github.com/brafales/piulades-bot/pinchito"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Crear struct {
	ChatID         int64
	Bot            *tgbotapi.BotAPI
	ActiveLogs     map[int]*pinchito.PlogData
	AuthToken      string
	PinchitoClient pinchito.Client
}

const keybReplySaveLog string = "SaveLog"
const keybReplyDiscardLog string = "DiscardLog"

const cmdNewLog string = "log"
const cmdEndLog string = "end_log"
const cmdStart string = "start"
const cmdCancel string = "cancel"

const titlePending string = "#Pending Title#"
const protagonistIdPending int = -1

func (c *Crear) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Crear")

	match, err := c.matchCommand(cmdStart, update.Message)
	if err != nil {
		return err
	}
	if match {
		return c.welcomeUser(update.Message)
	}

	match, err = c.matchCommand(cmdNewLog, update.Message)
	if err != nil {
		return err
	}
	if match {
		return c.startLogFromMessage(update.Message)
	}

	match, err = c.matchCommand(cmdEndLog, update.Message)
	if err != nil {
		return err
	}
	if match {
		return c.endLogFromMessage(update.Message)
	}

	match, err = c.matchCommand(cmdCancel, update.Message)
	if err != nil {
		return err
	}
	if match {
		_, err := c.reset(update.Message)
		return err
	}

	// If it's a CMD and it's not one of ours, this message is not for us
	// We are also allowing "//"
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() != "/" {
		return nil
	}

	isReply, err := c.hasInlineKeyboardReply(update)
	if err != nil {
		return err
	}
	if isReply {
		return c.handleKeyboardReply(update)
	}

	if c.userHasLogWithPendingTitle(update.Message) {
		return c.handleTitle(update.Message)
	}

	if c.userHasLogWithPendingProtagonist(update.Message) {
		return c.handleProtagonist(update.Message)
	}

	//This should be the last option; everything not handled before will be appended to the Log
	if c.userHasLogInProgress(update.Message) {
		return c.appendMessageToLog(update.Message)
	}

	//We shouldn't get here. Don't know what to do, send instructions
	tgMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "To Start a new log, type /"+cmdNewLog+" and then forward the messages you want to add. Once finished, type /"+cmdEndLog)
	c.Bot.Send(tgMessage)

	return nil
}

func (c *Crear) welcomeUser(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Hello "+message.From.FirstName+" "+message.From.LastName+"!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	c.Bot.Send(msg)

	c.sendMsg(message.Chat.ID, "To Start a new log, type /"+cmdNewLog+" and then forward the messages you want to add. Once finished, type /"+cmdEndLog)

	return nil
}

func (c *Crear) matchCommand(command string, message *tgbotapi.Message) (bool, error) {
	if message == nil {
		return false, nil
	}
	return regexp.MatchString("^"+command+"", message.Command())
}

func (c *Crear) hasInlineKeyboardReply(update tgbotapi.Update) (bool, error) {
	return update.CallbackQuery != nil, nil
}

func (c *Crear) reset(message *tgbotapi.Message) (bool, error) {
	delete(c.ActiveLogs, message.From.ID)

	tmp := tgbotapi.NewMessage(message.Chat.ID, "Operation cancelled. Back to square 1")
	tmp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	c.Bot.Send(tmp)

	return true, nil
}

func (c *Crear) startLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Starting a new Log for UserID ", message.From.ID)
	autor, err := c.PinchitoClient.GetUserFromTelegramUsername(message.From.UserName)
	if err != nil {
		c.sendMsg(message.Chat.ID, "Who are you? You don't seem to be a TruePinchito™. Talk with someone who can fix this")
		c.reset(message)
		return err
	}

	if c.userHasLogInProgress(message) {
		delete(c.ActiveLogs, message.From.ID)

		c.sendMsg(message.Chat.ID, "You already had a log in progress. Too bad. It has been discarded")
	}

	pinLog := pinchito.PlogData{}
	pinLog.Autor = autor.PinId
	c.ActiveLogs[message.From.ID] = &pinLog

	c.sendMsg(message.Chat.ID, "Ready to start creating a new Log.\nForward the messages you want to add. Once finished, type /"+cmdEndLog)

	return nil
}

func (c *Crear) endLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Ending existing log for UserID ", message.From.ID)

	pinLog := c.ActiveLogs[message.From.ID]
	if pinLog == nil {
		c.sendMsg(message.Chat.ID, "I haven't found any active log. Ignoring you")
		return nil
	}

	if len(pinLog.Text) == 0 {
		c.sendMsg(message.Chat.ID, "Why do you want to create an empty log? Try harder")
		return nil
	}
	pinLog.Data = time.Now().Unix()
	pinLog.Protagonista = protagonistIdPending

	c.askForTitle(message)

	return nil
}

func (c *Crear) appendMessageToLog(message *tgbotapi.Message) error {
	if message == nil {
		return nil
	}

	pinLog := c.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return nil
	}

	nick := c.GetNickFromMessage(message)
	pinLog.Text += pinmessage.BuildNewLogLine(nick, message)

	return nil
}

func (c *Crear) handleKeyboardReply(update tgbotapi.Update) error {
	oldMessage := update.CallbackQuery.Message
	text := ""

	if update.CallbackQuery.Data == keybReplySaveLog {
		logUrl, err := c.saveLog(update.CallbackQuery.From.ID)
		if err == nil {
			text = oldMessage.Text + "\n\nLog saved: " + logUrl
		} else {
			text = oldMessage.Text + "\n\nAn error occured while saving the log and it has been discarded:\n" + fmt.Sprint(err)
		}
	} else if update.CallbackQuery.Data == keybReplyDiscardLog {
		text = oldMessage.Text + "\nLog discarded"
	}

	editMessage := tgbotapi.NewEditMessageText(oldMessage.Chat.ID, oldMessage.MessageID, text)
	c.Bot.Send(editMessage)

	// If the upload succeeds, we need to delete it
	// If it fails, we delete it too and force the user to start over
	c.deleteLogFromUpdate(update)

	return nil
}

func (c *Crear) saveLog(userId int) (string, error) {

	pinLog := c.ActiveLogs[userId]

	uploadOp := &pinchito.JSONUploadOp{AuthToken: c.AuthToken, Upload: *pinLog}
	logUrl, err := c.PinchitoClient.UploadNewLog(uploadOp)

	return logUrl, err
}

func (c *Crear) deleteLogFromUpdate(update tgbotapi.Update) error {
	delete(c.ActiveLogs, update.CallbackQuery.From.ID)

	return nil
}

func (c *Crear) userHasLogInProgress(message *tgbotapi.Message) bool {
	return c.ActiveLogs[message.From.ID] != nil
}

func (c *Crear) userHasLogWithPendingTitle(message *tgbotapi.Message) bool {
	if message == nil {
		return false
	}

	pinLog := c.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return false
	}

	return pinLog.Titol == titlePending
}

func (c *Crear) askForTitle(message *tgbotapi.Message) error {
	//You can provide the title as an optional argument in cmdEndLog
	pinLog := c.ActiveLogs[message.From.ID]

	pinLog.Titol = titlePending
	c.sendMsg(message.Chat.ID, "Which title do you want the log to have?")

	return nil
}

func (c *Crear) handleTitle(message *tgbotapi.Message) error {
	pinLog := c.ActiveLogs[message.From.ID]
	pinLog.Titol = message.Text

	return c.askForProtagonist(message)
}

func (c *Crear) userHasLogWithPendingProtagonist(message *tgbotapi.Message) bool {
	if message == nil {
		return false
	}

	pinLog := c.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return false
	}

	return pinLog.Protagonista == protagonistIdPending
}

func (c *Crear) askForProtagonist(message *tgbotapi.Message) error {
	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Who is the protagonist of the Log?")

	users := c.PinchitoClient.GetPinchitoUsers()
	var buttons [][]tgbotapi.KeyboardButton
	for _, user := range users {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(user.PinNick)))
	}

	kbMarkup := tgbotapi.NewReplyKeyboard(buttons...)
	kbMarkup.OneTimeKeyboard = true
	kbMarkup.ResizeKeyboard = true
	tgMessage.ReplyMarkup = kbMarkup
	c.Bot.Send(tgMessage)

	return nil
}

func (c *Crear) handleProtagonist(message *tgbotapi.Message) error {

	user, err := c.PinchitoClient.GetUserFromPinchitoNick(message.Text)
	if err != nil {
		c.sendMsg(message.Chat.ID, "'"+message.Text+"' is not a TruePinchito™. Try again")
		return err
	}

	pinLog := c.ActiveLogs[message.From.ID]

	pinLog.Protagonista = user.PinId

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Using '"+user.PinNick+"' as your protagonist.")
	tgMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	c.Bot.Send(tgMessage)

	return c.sendLogSummary(message)
}

func (c *Crear) sendLogSummary(message *tgbotapi.Message) error {
	pinLog := c.ActiveLogs[message.From.ID]
	user, err := c.PinchitoClient.GetUserFromPinchitoNick(message.Text)
	if err != nil {
		c.sendMsg(message.Chat.ID, "'"+message.Text+"' is not a TruePinchito™. Try again")
		return err
	}

	c.sendMsg(message.Chat.ID, "I've created the following Log:")

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, pinLog.Titol+" (featuring "+user.PinNick+")\n\n"+pinLog.Text)
	inlineButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Save Log", keybReplySaveLog),
		tgbotapi.NewInlineKeyboardButtonData("Discard Log", keybReplyDiscardLog))
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineButtons)
	c.Bot.Send(tgMessage)

	return nil
}

func (c *Crear) sendMsg(chatID int64, text string) {
	tgMessage := tgbotapi.NewMessage(chatID, text)
	c.Bot.Send(tgMessage)
}

func (c *Crear) GetNickFromMessage(message *tgbotapi.Message) string {
	nick := ""
	if message.ForwardFrom != nil {
		user, err := c.PinchitoClient.GetUserFromTelegramUsername(message.ForwardFrom.UserName)
		if err == nil {
			nick = user.PinNick
		} else if len(message.ForwardFrom.UserName) > 0 {
			nick = message.ForwardFrom.UserName
		} else {
			nick = message.ForwardFrom.FirstName + " " + message.ForwardFrom.LastName
		}
	}

	return nick
}
