package handler

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"errors"
	pinmessage "github.com/brafales/piulades-bot/message"
	"github.com/brafales/piulades-bot/pinchito"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
)

type messageInfo struct {
	timestamp int
	message   tgbotapi.Message
}
type byTime []messageInfo

func (a byTime) Len() int           { return len(a) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTime) Less(i, j int) bool { return a[i].timestamp < a[j].timestamp }

type LogCreationData struct {
	logData  pinchito.PlogData
	messages []messageInfo
}

type Crear struct {
	ChatID         int64
	Bot            *tgbotapi.BotAPI
	ActiveLogs     map[int]*LogCreationData
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

	private, err := c.isPrivateMessage(update.Message)
	if err != nil || !private {
		//We are only dealing with private messages
		return err
	}

	match, err := c.matchCommand(cmdStart, update.Message)
	if err != nil {
		return err
	}
	if match {
		return c.welcomeUser(update.Message)
	}

	match, err = c.matchCommand(cmdNewLog, update.Message)
	if err != nil || (!match && update.Message != nil && !c.userHasLogInProgress(update.Message)) {
		//if we receive a message, it's not to start a new log and we have no log in progress, ignore it too
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
		log.Print("Unknown CMD:" + update.Message.Command())
		c.sendMsg(update.Message.Chat.ID, "I don't know what you mean with '/"+update.Message.Command()+"' Check the list of commands by typing '/' and disable your keyboard's auto-correct system")
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
		return c.appendMessageToLog(update)
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

	pinLog := LogCreationData{}
	pinLog.logData.Autor = autor.PinId
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

	if len(pinLog.messages) == 0 {
		c.sendMsg(message.Chat.ID, "Why do you want to create an empty log? Try harder")
		return nil
	}
	pinLog.logData.Data = time.Now().Unix()
	pinLog.logData.Protagonista = protagonistIdPending

	c.askForTitle(message)

	return nil
}

func (c *Crear) appendMessageToLog(update tgbotapi.Update) error {
	message := update.Message
	if message == nil {
		return nil
	}

	pinLog := c.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return nil
	}

	pinLog.messages = append(pinLog.messages, messageInfo{update.UpdateID, *message})

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

	if pinLog == nil {
		return "", errors.New("User had no log in progress")
	}

	uploadOp := &pinchito.JSONUploadOp{AuthToken: c.AuthToken, Upload: pinLog.logData}
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

	return pinLog.logData.Titol == titlePending
}

func (c *Crear) askForTitle(message *tgbotapi.Message) error {
	pinLog := c.ActiveLogs[message.From.ID]

	pinLog.logData.Titol = titlePending
	c.sendMsg(message.Chat.ID, "Which title do you want the log to have?")

	return nil
}

func (c *Crear) handleTitle(message *tgbotapi.Message) error {
	pinLog := c.ActiveLogs[message.From.ID]
	pinLog.logData.Titol = message.Text

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

	return pinLog.logData.Protagonista == protagonistIdPending
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

	pinLog.logData.Protagonista = user.PinId

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Using '"+user.PinNick+"' as your protagonist.")
	tgMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	c.Bot.Send(tgMessage)

	c.sortMesagesAndCreateText(pinLog)

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

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, pinLog.logData.Titol+" (featuring "+user.PinNick+")\n\n"+pinLog.logData.Text)
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
func (c *Crear) isPrivateMessage(message *tgbotapi.Message) (bool, error) {
	return message == nil || message.Chat.IsPrivate(), nil
}
func (c *Crear) sortMesagesAndCreateText(logCreationData *LogCreationData) {
	sort.Sort(byTime(logCreationData.messages))

	for _, messageInfo := range logCreationData.messages {
		nick := c.GetNickFromMessage(&messageInfo.message)
		logCreationData.logData.Text += pinmessage.BuildNewLogLine(nick, &messageInfo.message)

	}
}
