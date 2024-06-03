package utils

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PrintWelcome() {
	log.Printf("Lemonoid awakens!")
}

func PrintMsg(update tgbotapi.Update) {
	log.Printf("[%s-%d] [MsgId: %d]: %s",
		update.Message.From.UserName,
		update.Message.From.ID,
		update.Message.MessageID,
		update.Message.Text)
}
