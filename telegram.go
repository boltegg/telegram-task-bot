package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func ProcessTelegramBot() {

	t := NewTelegram()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := t.bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			go t.MessageHandler(&update)
		} else if update.CallbackQuery != nil {
			go t.CallbackHandler(update.CallbackQuery)
		}

	}
}

func NewTelegram() *Telegram {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotApiKey)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Authorized on account %s", bot.Self.UserName)

	return &Telegram{bot}
}

func (t *Telegram) CallbackHandler(c *tgbotapi.CallbackQuery) {

	//logrus.Infof("[%d] %s", c.Message.Chat.ID, c.Data)

	cmd := strings.Split(c.Data, " ")

	// close task
	if len(cmd) == 2 && cmd[0] == "close" {
		i, err := strconv.ParseInt(cmd[1], 10, 64)
		if err != nil {
			t.CallbackAnswer(c.ID, "Error: "+err.Error())
			return
		}
		err = CloseTask(c.Message.Chat.ID, i)
		if err != nil {
			t.CallbackAnswer(c.ID, "Error: "+err.Error())
			return
		}

		t.CallbackAnswer(c.ID, "Task Closed!")

		editMessage, err := MsgEditReplyMarkup(c.Message, i)
		if err != nil {
			t.SendError(c.Message.Chat.ID, err)
			return
		}

		t.bot.Send(editMessage)

	}

	// open task
	if len(cmd) == 2 && cmd[0] == "open" {
		i, err := strconv.ParseInt(cmd[1], 10, 64)
		if err != nil {
			t.CallbackAnswer(c.ID, "Error: "+err.Error())
			return
		}
		err = OpenTask(c.Message.Chat.ID, i)
		if err != nil {
			t.CallbackAnswer(c.ID, "Error: "+err.Error())
			return
		}

		t.CallbackAnswer(c.ID, "Task Opened!")

		editMessage, err := MsgEditReplyMarkup(c.Message, i)
		if err != nil {
			t.SendError(c.Message.Chat.ID, err)
			return
		}

		t.bot.Send(editMessage)
	}

}

func (t *Telegram) SendError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, "Error: "+err.Error())
	t.bot.Send(msg)
}

func (t *Telegram) CallbackAnswer(callbackQueryID string, text string) {
	ccc := tgbotapi.CallbackConfig{CallbackQueryID: callbackQueryID, Text: text}
	_, _ = t.bot.AnswerCallbackQuery(ccc)
}

func (t *Telegram) MessageHandler(update *tgbotapi.Update) {

	var m = update.Message

	//logrus.Infof("[%d] %s", m.Chat.ID, m.Text)

	if !m.IsCommand() {
		// create task
		i, err := NewTask(m.Chat.ID, m.Chat.ID, m.Text).Create()
		if err != nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, "Task creation error")
			t.bot.Send(msg)
		}
		msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("Task created /%d", i))
		t.bot.Send(msg)
		return
	}

	i, err := strconv.ParseInt(m.Command(), 10, 64)
	if err == nil {
		// print task
		//m := tgbotapi.NewMessage()
		msg, err := MsgTask(m, i)
		if err != nil {
			_, err = t.bot.Send(tgbotapi.NewMessage(m.Chat.ID, "Error: "+err.Error()))
		}

		_, err = t.bot.Send(msg)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	switch m.Command() {
	case "today":
		//TODO: print tasklist today
	case "tommorow":
		//TODO: print tasklist tommorow
	case "all":
		msg := PrintTaskList(m)
		t.bot.Send(msg)
		return
	case "setdue":
		//TODO: set due date for task

		//TODO: more commands
	}

	//logrus.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//msg.ReplyToMessageID = update.Message.MessageID
	//
	//t.bot.Send(msg)
}

// example use buttons
// msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard:[][]tgbotapi.KeyboardButton{[]tgbotapi.KeyboardButton{{Text:"dfdf"}}}}

func MsgTask(m *tgbotapi.Message, i int64) (msg tgbotapi.MessageConfig, err error) {

	t, err := GetTask(m.Chat.ID, int64(i))
	if err != nil {
		return
	}

	msg = tgbotapi.NewMessage(m.Chat.ID, t.PrintMessage())
	msg.ReplyMarkup = t.PrintKeyboard()
	msg.ParseMode = "HTML"

	return
}

func MsgEditReplyMarkup(m *tgbotapi.Message, i int64) (msg tgbotapi.EditMessageReplyMarkupConfig, err error) {

	t, err := GetTask(m.Chat.ID, int64(i))
	if err != nil {
		return
	}

	msg = tgbotapi.NewEditMessageReplyMarkup(m.Chat.ID, m.MessageID, t.PrintKeyboard())
	return
}

func PrintTaskList(m *tgbotapi.Message) (msg tgbotapi.MessageConfig) {

	tasks, err := GetAllTask(m.Chat.ID)
	if err != nil {
		msg = tgbotapi.NewMessage(m.Chat.ID, err.Error())
		return
	}

	message := "List of all your tasks:\n"
	for _, task := range tasks {
		message += fmt.Sprintf("/%d <b>%s</b>\n", task.TaskID, task.Subject)
	}
	msg = tgbotapi.NewMessage(m.Chat.ID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	return
}

func (t *Task) PrintMessage() string {
	return fmt.Sprintf("<b>Task #%d: %s</b>\n\n%s", t.TaskID, t.Subject, t.Description)
}

func (t *Task) PrintKeyboard() tgbotapi.InlineKeyboardMarkup {

	var button_status tgbotapi.InlineKeyboardButton

	if t.Status == "open" {
		button_status = btn("Close", fmt.Sprintf("close %d", t.TaskID))
	} else if t.Status == "close" {
		button_status = btn("Open", fmt.Sprintf("open %d", t.TaskID))
	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{button_status},
		},
	}
}

func btn(text, callbackData string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callbackData}
}
