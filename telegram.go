package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"

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
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		t.RequestHandler(&update)
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
func (t *Telegram) RequestHandler(update *tgbotapi.Update) {

	logrus.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	t.bot.Send(msg)
}


// example use buttons
// msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard:[][]tgbotapi.KeyboardButton{[]tgbotapi.KeyboardButton{{Text:"dfdf"}}}}