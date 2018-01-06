package handler

import (
	"log"
	"regexp"
	"time"
	"strconv"
	"fmt"
	pinmessage "github.com/brafales/piulades-bot/message"

	"github.com/brafales/piulades-bot/pinchito"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

type Crear struct {
	ChatID     int64
	Bot        *tgbotapi.BotAPI
	ActiveLogs map[int]*pinchito.PlogData
	AuthToken  string
}

const keybReplySaveLog string = "SaveLog"
const keybReplyDiscardLog string = "DiscardLog"

const cmdNewLog string = "log"
const cmdEndLog string = "end_log"
const cmdStart string = "start"
const cmdCancel string = "cancel"

const titlePending string = "#Pending Title#"
const protagonistIdPending int = -1


func (t *Crear) Handle(update tgbotapi.Update) error {
	log.Println("Handling with Crear")

	match, err := t.matchCommand(cmdStart, update.Message)
	if err != nil {
		return err
	}
	if match {
		return t.welcomeUser(update.Message)
	}

	match, err = t.matchCommand(cmdNewLog, update.Message)
	if err != nil {
		return err
	}
	if match {
		return t.startLogFromMessage(update.Message)
	}


	match, err = t.matchCommand(cmdEndLog, update.Message)
	if err != nil {
		return err
	}
	if match {
		return t.endLogFromMessage(update.Message)
	}


	match, err = t.matchCommand(cmdCancel, update.Message)
	if err != nil {
		return err
	}
	if match {
		_,err := t.reset(update.Message)
		return err
	}

	// If it's a CMD and it's not one of ours, this message is not for us
	// We are also allowing "//"
	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() != "/"{
		log.Print("Unknown CMD:" + update.Message.Command())
		t.sendMsg(update.Message.Chat.ID, "I don't know what you mean with '/" + update.Message.Command() + "' Check the list of commands by typing '/' and disable your keyboard's auto-correct system")
		return nil
	}





	isReply, err := t.hasInlineKeyboardReply(update)
	if err != nil {
		return err
	}
	if isReply {
		return t.handleKeyboardReply(update)
	}


	if t.userHasLogWithPendingTitle(update.Message) {
		return t.handleTitle(update.Message)
	}


	if t.userHasLogWithPendingProtagonist(update.Message) {
		return t.handleProtagonist(update.Message)
	}


	//This should be the last option; everything not handled before will be appended to the Log
	if t.userHasLogInProgress(update.Message) {
		return t.appendMessageToLog(update.Message)
	}


	//We shouldn't get here. Don't know what to do, send instructions
	tgMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "To Start a new log, type /" +cmdNewLog+ " and then forward the messages you want to add. Once finished, type /" + cmdEndLog)
	t.Bot.Send(tgMessage)

	return nil
}


func (t *Crear) welcomeUser(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Hello "+ message.From.FirstName + " " + message.From.LastName + "!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	t.Bot.Send(msg)

	t.sendMsg(message.Chat.ID, "To Start a new log, type /" +cmdNewLog+ " and then forward the messages you want to add. Once finished, type /" + cmdEndLog)

	return nil
}

func (t *Crear) matchCommand(command string, message *tgbotapi.Message) (bool, error) {
	if message == nil {
		return false, nil
	}
	return regexp.MatchString("^" + command + "", message.Command())
}

func (t *Crear) hasInlineKeyboardReply(update tgbotapi.Update) (bool, error) {
	return update.CallbackQuery != nil, nil
}

func (t *Crear) reset(message *tgbotapi.Message) (bool, error) {
	delete(t.ActiveLogs, message.From.ID)

	tmp := tgbotapi.NewMessage(message.Chat.ID, "Operation cancelled. Back to square 1")
	tmp.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	t.Bot.Send(tmp)

	return true, nil
}


func (t *Crear) startLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Starting a new Log for UserID ", message.From.ID)

	if t.userHasLogInProgress(message) {
		delete(t.ActiveLogs, message.From.ID)

		t.sendMsg(message.Chat.ID, "You already had a log in progress. Too bad. It has been discarded")
	}

	t.ActiveLogs[message.From.ID] = &pinchito.PlogData{}

	t.sendMsg(message.Chat.ID, "Ready to start creating a new Log.\nForward the messages you want to add. Once finished, type /" + cmdEndLog)

	return nil
}

func (t *Crear) endLogFromMessage(message *tgbotapi.Message) error {
	log.Println("Ending existing log for UserID ", message.From.ID)

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		t.sendMsg(message.Chat.ID, "I haven't found any active log. Ignoring you")
		return nil
	}

	if len(pinLog.Text) == 0 {
		t.sendMsg(message.Chat.ID, "Why do you want to create an empty log? Try harder")
		return nil
	}

	autor, err := pinchito.GetUserFromTelegramUsername(message.From.UserName)
	if err != nil {
		return err
	}
	pinLog.Autor = autor.PinId
	pinLog.Data = time.Now().Unix()
	pinLog.Protagonista = protagonistIdPending

	t.askForTitle(message)

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

	line := pinmessage.BuildNewLogLine(message)
	if len(line) > 0 {
		pinLog.Text += line
	}
	// TODO Handle other types of messages
	// (Audio, img, ...)

	return nil
}

func (t *Crear) handleKeyboardReply(update tgbotapi.Update) error {
	host := os.Getenv("PINCHITO_HOST")

	oldMessage := update.CallbackQuery.Message
	text := ""
	if update.CallbackQuery.Data == keybReplySaveLog {
		logId, err := t.saveLog(update.CallbackQuery.From.ID)
		if err == nil {
			text = oldMessage.Text + "\n\nLog saved: http://" + host + "/" + strconv.Itoa(logId)
		} else {
			text = oldMessage.Text + "\n\nAn error occured while saving the log and it has been discarded:\n" + fmt.Sprint(err)
		}
	} else if update.CallbackQuery.Data == keybReplyDiscardLog {
		text = oldMessage.Text + "\nLog discarded"
	}

	editMessage := tgbotapi.NewEditMessageText(oldMessage.Chat.ID, oldMessage.MessageID, text)
	t.Bot.Send(editMessage)

	// If the upload succeeds, we need to delete it
	// If it fails, we delete it too and force the user to start over
	t.deleteLogFromUpdate(update)

	return nil
}

func (t *Crear) saveLog(userId int) (int, error) {

	pinLog := t.ActiveLogs[userId]

	uploadOp := &pinchito.JSONUploadOp{AuthToken:t.AuthToken, Upload:*pinLog}
	logId, err := pinchito.UploadNewLog(uploadOp)

	return logId, err
}

func (t *Crear) deleteLogFromUpdate(update tgbotapi.Update) error {
	delete(t.ActiveLogs, update.CallbackQuery.From.ID)

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

func (t *Crear) askForTitle(message *tgbotapi.Message) error {
	//You can provide the title as an optional argument in cmdEndLog
	pinLog := t.ActiveLogs[message.From.ID]

	pinLog.Titol = titlePending
	t.sendMsg(message.Chat.ID, "Which title do you want the log to have?")

	return nil
}

func (t *Crear) handleTitle(message *tgbotapi.Message) error {
	pinLog := t.ActiveLogs[message.From.ID]
	pinLog.Titol = message.Text

	return t.askForProtagonist(message)
}

func (t *Crear) userHasLogWithPendingProtagonist(message *tgbotapi.Message) bool {
	if message == nil {
		return false
	}

	pinLog := t.ActiveLogs[message.From.ID]
	if pinLog == nil {
		return false
	}

	return pinLog.Protagonista == protagonistIdPending
}

func (t *Crear) askForProtagonist(message *tgbotapi.Message) error {
	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Who is the protagonist of the Log?")

	users := pinchito.GetPinchitoUsers()
	var buttons [][]tgbotapi.KeyboardButton
	for _, user := range users {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(user.PinNick)))
	}

	kbMarkup := tgbotapi.NewReplyKeyboard(buttons...)
	kbMarkup.OneTimeKeyboard = true
	kbMarkup.ResizeKeyboard = true
	tgMessage.ReplyMarkup = kbMarkup
	t.Bot.Send(tgMessage)

	return nil
}

func (t *Crear) handleProtagonist(message *tgbotapi.Message) error {

	user,err := pinchito.GetUserFromPinchitoNick(message.Text)
	if err != nil {
		t.sendMsg(message.Chat.ID, "'" + message.Text + "' is not a TruePinchito™. Try again")
		return err
	}

	pinLog := t.ActiveLogs[message.From.ID]

	pinLog.Protagonista = user.PinId

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, "Using '" + user.PinNick + "' as your protagonist.")
	tgMessage.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	t.Bot.Send(tgMessage)

	return t.sendLogSummary(message)
}

func (t *Crear) sendLogSummary(message *tgbotapi.Message) error {
	pinLog := t.ActiveLogs[message.From.ID]
	user,err := pinchito.GetUserFromPinchitoNick(message.Text)
	if err != nil {
		t.sendMsg(message.Chat.ID, "'" + message.Text + "' is not a TruePinchito™. Try again")
		return err
	}

	t.sendMsg(message.Chat.ID, "I've created the following Log:")

	tgMessage := tgbotapi.NewMessage(message.Chat.ID, pinLog.Titol + " (featuring " + user.PinNick + ")\n\n" + pinLog.Text)
	inlineButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Save Log", keybReplySaveLog),
		tgbotapi.NewInlineKeyboardButtonData("Discard Log", keybReplyDiscardLog))
	tgMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineButtons)
	t.Bot.Send(tgMessage)

	return nil
}

func (t *Crear) sendMsg(chatID int64, text string) {
	tgMessage := tgbotapi.NewMessage(chatID, text)
	t.Bot.Send(tgMessage)
}