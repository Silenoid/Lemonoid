package telegram

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Silenoid/Lemonoid/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var token string
var tgClient *bot.Bot
var bgCtx context.Context
var startTime time.Time

var CHATID_LORD int64 = 449697032
var CHATID_CHICECE int64 = -1001623264158

var PROMPT_THEMES []string = []string{
	"una avventura fantastica",
	"una lettera scritta in epoca vittoriana",
	"una storia di Natale",
	"un vecchio articolo di giornale",
	"una missiva in epoca medievale",
	"un antico editto romano",
	"una storia dell'orrore",
	"un discorso politico",
	"un sermone fatto in chiesa",
	"un romanzo rosa",
	"una rivista di gossip",
	"un canto della divina commedia",
	"una poesia ermetica",
	"un libro per bambini",
	"il cronista di una partita di calcio",
	"un elogio funebre ad un funerale",
	"uno spot pubblicitario",
	"un riassunto delle puntate precedenti",
	"un triste ed angosciante racconto sovietico",
	"un passo della Bibbia",
	"un discorso fra Gen Z in un linguaggio dank contenente frequenti riferimenti a meme",
}

func Initialize() {
	token = utils.TokenTelegram

	customOptions := []bot.Option{
		bot.WithMessageTextHandler("/tldrpro", bot.MatchTypePrefix, handlerTldrPro),
		bot.WithMessageTextHandler("/tldr", bot.MatchTypePrefix, handlerTldr),
		bot.WithMessageTextHandler("/help", bot.MatchTypePrefix, handlerHelp),
		bot.WithMessageTextHandler("/status", bot.MatchTypePrefix, handlerStatus),
		bot.WithMessageTextHandler("/stamoce", bot.MatchTypePrefix, handlerStamoce),
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCheckInitTimeout(60 * time.Second),
		// bot.WithDebug(),
	}

	telegramBot, err := bot.New(token, customOptions...)
	if err != nil {
		log.Panic("Error during bot initialization -> ", err)
	}

	tgClient = telegramBot

	startTime = time.Now()
}

func Listen() {
	// Gracefull shotdown
	backgroundContext, interruptCallback := signal.NotifyContext(context.Background(), os.Interrupt)
	defer interruptCallback()

	bgCtx = backgroundContext

	tgClient.SetMyCommands(backgroundContext, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "help", Description: "s'aiutamo"},
			{Command: "status", Description: "'ndo stamo"},
			{Command: "tldrpro", Description: "che se dice, co Graizano Maronne"},
			{Command: "tldr", Description: "che se dice"},
			{Command: "stamoce", Description: "quanno se vedemo"},
		},
	})

	log.Println("[Telegram client] Telegram client is listening...")

	SendMessage(CHATID_LORD, "Lemonoid awakened at "+utils.ToReadableDate(startTime))

	tgClient.Start(bgCtx)
}

func isAllowedChatId(chatId int64) bool {
	return (chatId == CHATID_CHICECE || chatId == CHATID_LORD)
}

func sendAudio(update *models.Update, audioPath string) {
	audioFileContent, err := os.ReadFile(audioPath)

	if err != nil {
		log.Println("Errore nella lettura del file audio " + audioPath)
		return
	}

	sendVoiceParams := bot.SendVoiceParams{
		ChatID:          update.Message.Chat.ID,
		MessageThreadID: update.Message.MessageThreadID,
		Voice: &models.InputFileUpload{
			Filename: filepath.Base(audioPath),
			Data:     bytes.NewReader(audioFileContent),
		},
	}

	tgClient.SendVoice(bgCtx, &sendVoiceParams)
}

func RespondWithText(update *models.Update, text string) {
	params := &bot.SendMessageParams{
		ChatID:          update.Message.Chat.ID,
		MessageThreadID: update.Message.MessageThreadID,
		Text:            text,
	}

	tgClient.SendMessage(bgCtx, params)
}

func SendMessage(chatId int64, text string) {
	params := &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	}

	tgClient.SendMessage(bgCtx, params)
}
