package utils

import (
	"log"

	"github.com/go-telegram/bot/models"
)

func PrintWelcome() {
	log.Printf("Lemonoid awakens!")
}

func PrintMsg(update *models.Update) {
	log.Printf("[%s-%d] [MsgId: %d]: %s",
		update.Message.From.Username,
		update.Message.From.ID,
		update.Message.ID,
		update.Message.Text)
}
