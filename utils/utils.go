package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

func SendMessage(message string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func GetUnixtimestamp(offset int64) int64 {
	return time.Now().Unix() + offset
}
