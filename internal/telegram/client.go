package telegram

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/gemini"
	"github.com/Silenoid/Lemonoid/internal/history"
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

const TLDR_WAIT_TIME = time.Hour * 24

var lastTLDRTime time.Time
var isFirstTLDR = false

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
			{Command: "tldr", Description: "che se dice"},
			{Command: "stamoce", Description: "quanno se vedemo"},
			{Command: "sperimentale", Description: "non clickare mai"},
		},
	})

	log.Println("[Telegram client] Telegram client is listening...")

	SendMessage(CHATID_LORD, "Lemonoid awakened at "+utils.ToReadableDate(startTime))

	tgClient.Start(bgCtx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// if update.Message != nil {
	// 	b.SendMessage(ctx, &bot.SendMessageParams{
	// 		ChatID: update.Message.Chat.ID,
	// 		Text:   update.Message.Text,
	// 	})
	// }

	// params := &bot.SendPhotoParams{
	// 	ChatID:  update.Message.Chat.ID,
	// 	Photo:   &models.InputFileString{Data: "AgACAgIAAxkDAAIBOWJimnCJHQJiJ4P3aasQCPNyo6mlAALDuzEbcD0YSxzjB-vmkZ6BAQADAgADbQADJAQ"},
	// 	Caption: "Preloaded Facebook logo",
	// }

	// b.SendPhoto(ctx, params)

	if update.Message != nil {
		utils.PrintMsg(update)

		// update.Message.Date
		// time.unix
		// if update.Message.Time().Before(startTime) {
		// 	log.Printf("Not processing message %d from %s-%d because of previous messages cleanup",
		// 		update.Message.From.ID,
		// 		update.Message.From.UserName,
		// 		update.Message.From.ID)

		// }

		if isAllowedChatId(update.Message.Chat.ID) {
			err := processAndDispatch(update)

			if err != nil {
				log.Printf("Error during processing of message from %s-%d with ID %d\n%v",
					update.Message.From.Username,
					update.Message.From.ID,
					update.Message.ID,
					err)
			}
		} else {
			log.Printf("Not processing message %d from %s-%d because of chat exclusions",
				update.Message.From.ID,
				update.Message.From.Username,
				update.Message.From.ID)
		}
	}
}

func isAllowedChatId(chatId int64) bool {
	return (chatId == CHATID_CHICECE || chatId == CHATID_LORD)
}

func processAndDispatch(update *models.Update) error {
	message := update.Message
	chat := update.Message.Chat
	senderName := update.Message.From.Username
	senderId := update.Message.From.ID
	textReceived := update.Message.Text

	if textReceived != "" {
		history.AddMessageToChatHistory(chat.ID, senderId, senderName, message.Text)
	}

	log.Println("Processing message: " + textReceived)

	switch textReceived {
	case "/help":
		return handlerHelp(update)
	case "/status":
		return handlerStatus(update)
	case "/tldr":
		return handlerTldr(update)
	case "/stamoce":
		return handlerStamoce(update)
	default:
		log.Printf("Command '%s' has not been recognized. Nope.", textReceived)
	}
	return nil
}

func handlerHelp(update *models.Update) error {
	SendMessage(update.Message.Chat.ID, `Aò a manzo, eccote du seppie e ttre ppiovre de aiuto:
	/help	l'hai usato mò a cojone, ma che sei frocio?
	/tldr	azzì questo teggenera er tuloddonrì
	/status	je chiedi mammamiacommestaaa
	`)
	return nil
}

func handlerStatus(update *models.Update) error {
	ElevenLabsSubStatus := elevenlabs.GetSubscriptionStatus()
	SendMessage(update.Message.Chat.ID, ElevenLabsSubStatus)
	return nil
}

func handlerTldr(update *models.Update) error {
	if !isFirstTLDR && time.Now().Before(lastTLDRTime.Add(TLDR_WAIT_TIME)) {
		SendMessage(update.Message.Chat.ID, "A cojò, li mortacci stracci tua ma che me voj rovinà? So ppasati solo "+utils.ToReadableSince(time.Now(), lastTLDRTime)+" da nantro vocale. Statte bbono pe' n'antri "+utils.ToReadableHowLongTo(time.Now(), lastTLDRTime, TLDR_WAIT_TIME))
		return nil
	}

	chatHistory := history.GetChatHistory(update.Message.Chat.ID)
	pickedPromptTheme := utils.PickFromArray(PROMPT_THEMES)
	// pickedVoice := utils.PickFromArray(elevenlabs.BASIC_VOICES)
	// Impose G Maronne
	pickedVoice := "ml5JfpB48j688Rpbbz2M"

	llmPromptBuilder := strings.Builder{}
	llmPromptBuilder.WriteString("Genera un riassunto della seguente chat come se fosse ")
	llmPromptBuilder.WriteString(pickedPromptTheme)

	switch pickedVoice {
	case "IzoLtTXseyrunESwWmw3": // Se è M TODO: definisci enum o tipo
		llmPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'devastaaaante' e facendo paragoni col Giappone:\n")
	case "i86lB8eIKMQcO470EIFz", "ml5JfpB48j688Rpbbz2M": // // Se è G o G Maronne
		llmPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'WAGOOOOO' ed concludendo, alla fine, suggerendo un piatto di pasta insolito da cucinare:\n")
	case "d9Gr3L3YR4d9Sf9Gt8cV": // Se è S
		llmPromptBuilder.WriteString(", utilizzando almeno una volta ciascuno i termini 'non ironicamente', 'cringe' e 'è tutta colpa di Enzo':\n")
	default:
		llmPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'grottesco':\n")
	}

	llmPromptBuilder.WriteString(chatHistory)
	llmPrompt := llmPromptBuilder.String()

	log.Printf("Prompt a tema '%s' con voce '%s': %s", pickedPromptTheme, pickedVoice, llmPrompt)
	SendMessage(update.Message.Chat.ID, "Tema utilizzato per il prompt: "+pickedPromptTheme)

	generatedStory, err := gemini.GenerateStory(llmPrompt)
	if err != nil {
		if strings.Contains(update.Message.Text, "exceeded your current quota") {
			SendMessage(update.Message.Chat.ID, "Ao, so ffiniti li sordi pe generà er testo li mortacci stracci")
		}
		SendMessage(CHATID_LORD, err.Error())
		return err
	}

	lastTLDRTime = time.Now()
	isFirstTLDR = false

	generatedAudioPath, err := elevenlabs.GenerateVoiceNarration(generatedStory, pickedVoice)
	if err != nil {
		SendMessage(update.Message.Chat.ID, "Errore nella generazione vocale, dunque beccate solo er testo generato e muto:\n"+generatedStory)
		return err
	}

	sendAudio(update.Message.Chat.ID, generatedAudioPath)
	SendMessage(CHATID_LORD, "Generated story:\n"+generatedStory)
	return nil
}

func handlerStamoce(update *models.Update) error {
	// Telegram supports only up to 10 options
	pollOptions := []models.InputPollOption{
		{Text: "Lunneddì sera"},
		{Text: "Marteddì sera"},
		{Text: "Mercoleddì sera"},
		{Text: "Gioveddì sera"},
		{Text: "Venerdì sera"},
		{Text: "Sabato mattina"},
		{Text: "Sabato sera"},
		{Text: "Domenica mattina"},
		{Text: "Domenica sera"},
	}

	tgClient.SendPoll(
		bgCtx,
		&bot.SendPollParams{
			ChatID:                update.Message.Chat.ID,
			Question:              "Quanno ce stamo?",
			Options:               pollOptions,
			IsAnonymous:           bot.False(),
			AllowsMultipleAnswers: *bot.True(),
			Explanation:           "t'o devo pure spiegà?",
		},
	)

	return nil
}

func sendAudio(chatId int64, audioPath string) {
	audioFileContent, err := os.ReadFile(audioPath)

	if err != nil {
		log.Println("Errore nella lettura del file audio " + audioPath)
		return
	}

	sendVoiceParams := bot.SendVoiceParams{
		ChatID: chatId,
		Voice: &models.InputFileUpload{
			Filename: filepath.Base(audioPath),
			Data:     bytes.NewReader(audioFileContent),
		},
	}

	tgClient.SendVoice(bgCtx, &sendVoiceParams)
}

func SendMessage(chatId int64, text string) {
	params := &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	}

	tgClient.SendMessage(bgCtx, params)
}
