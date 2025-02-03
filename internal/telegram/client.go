package telegram

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/history"
	"github.com/Silenoid/Lemonoid/internal/openai"
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
		bot.WithDefaultHandler(defaultHandler),
		// bot.WithDebug(),
	}

	telegramBot, err := bot.New(token, customOptions...)
	if err != nil {
		log.Panic("Error during bot initialization -> ", err)
	}

	backgroundContext, interruptCallback := signal.NotifyContext(context.Background(), os.Interrupt)
	defer interruptCallback()

	bgCtx = backgroundContext
	tgClient = telegramBot

	tgClient.Start(bgCtx)

	startTime = time.Now()
	SendMessage(CHATID_LORD, "Lemonoid awakened at "+utils.ToReadableDate(startTime))
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}

	params := &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileString{Data: "AgACAgIAAxkDAAIBOWJimnCJHQJiJ4P3aasQCPNyo6mlAALDuzEbcD0YSxzjB-vmkZ6BAQADAgADbQADJAQ"},
		Caption: "Preloaded Facebook logo",
	}

	b.SendPhoto(ctx, params)
}

func Listen() {
	if len(token) == 0 {
		panic("Telegram token is not set")
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updatesChannel := tgClient.GetUpdatesChan(updateConfig)
	updatesChannel.Clear()

	log.Println("[Telegram client] Telegram client is listening...")

	for update := range updatesChannel {

		if update.Message != nil {
			utils.PrintMsg(update)

			if update.Message.Time().Before(startTime) {
				log.Printf("Not processing message %d from %s-%d because of previous messages cleanup",
					update.Message.From.ID,
					update.Message.From.UserName,
					update.Message.From.ID)

			}

			if isAllowedChatId(update.Message.Chat.ID) {
				err := processAndDispatch(update)

				if err != nil {
					log.Printf("Error during processing of message from %s-%d with ID %d\n%v",
						update.Message.From.UserName,
						update.Message.From.ID,
						update.Message.MessageID,
						err)
				}
			} else {
				log.Printf("Not processing message %d from %s-%d because of chat exclusions",
					update.Message.From.ID,
					update.Message.From.UserName,
					update.Message.From.ID)
			}

		}
	}
}

func isAllowedChatId(chatId int64) bool {
	return (chatId == CHATID_CHICECE || chatId == CHATID_LORD)
}

func processAndDispatch(update tgbotapi.Update) error {
	message := update.Message
	chat := update.Message.Chat
	senderName := update.Message.From.UserName
	senderId := update.Message.From.ID
	textReceived := update.Message.Text

	if textReceived != "" {
		history.AddMessageToChatHistory(chat.ID, senderId, senderName, message.Text)
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "help":
			return handlerHelp(update)
		case "status":
			return handlerStatus(update)
		case "tldr":
			return handlerTldr(update)
		case "stamoce":
			return handlerStamoce(update)
		default:
			log.Printf("Command %s has not been recognized. Nope.", update.Message.Command())
		}
	}
	return nil
}

func handlerHelp(update tgbotapi.Update) error {
	SendMessage(update.Message.Chat.ID, `Aò a manzo, eccote du seppie e ttre ppiovre de aiuto:
	/help	l'hai usato mò a cojone, ma che sei frocio?
	/tldr	azzì questo teggenera er tuloddonrì
	/status	je chiedi mammamiacommestaaa
	`)
	return nil
}

func handlerStatus(update tgbotapi.Update) error {
	ElevenLabsSubStatus := elevenlabs.GetSubscriptionStatus()
	// TODO: get OpenAI usage with a request (see openai client.go)
	SendMessage(update.Message.Chat.ID, ElevenLabsSubStatus)
	return nil
}

func handlerTldr(update tgbotapi.Update) error {
	chatHistory := history.GetChatHistory(update.Message.Chat.ID)
	pickedPromptTheme := utils.PickFromArray(PROMPT_THEMES)
	// pickedVoice := utils.PickFromArray(elevenlabs.BASIC_VOICES)
	// Impose G Maronne
	pickedVoice := "ml5JfpB48j688Rpbbz2M"

	openAiPromptBuilder := strings.Builder{}
	openAiPromptBuilder.WriteString("Genera un riassunto della seguente chat come se fosse ")
	openAiPromptBuilder.WriteString(pickedPromptTheme)

	switch pickedVoice {
	case "IzoLtTXseyrunESwWmw3": // Se è M TODO: definisci enum o tipo
		openAiPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'devastaaaante' e facendo paragoni col Giappone:\n")
	case "i86lB8eIKMQcO470EIFz", "ml5JfpB48j688Rpbbz2M": // // Se è G o G Maronne
		openAiPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'WAGOOOOO' ed concludendo, alla fine, suggerendo un piatto di pasta insolito da cucinare:\n")
	case "d9Gr3L3YR4d9Sf9Gt8cV": // Se è S
		openAiPromptBuilder.WriteString(", utilizzando almeno una volta ciascuno i termini 'non ironicamente', 'cringe' e 'è tutta colpa di Enzo':\n")
	default:
		openAiPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'grottesco':\n")
	}

	openAiPromptBuilder.WriteString(chatHistory)
	openAiPrompt := openAiPromptBuilder.String()

	log.Printf("Prompt a tema '%s' con voce '%s': %s", pickedPromptTheme, pickedVoice, openAiPrompt)
	elevenLabsPrompt, err := openai.GenerateStory(openAiPrompt)
	if err != nil {
		return err
	}

	generatedAudioPath, err := elevenlabs.GenerateVoiceNarration(elevenLabsPrompt, pickedVoice)
	if err != nil {
		SendMessage(update.Message.Chat.ID, "Errore durante la generazione vocale: "+err.Error())
		return err
	}

	SendMessage(update.Message.Chat.ID, "Tema utilizzato per il prompt: "+pickedPromptTheme)
	sendAudio(update.Message.Chat.ID, generatedAudioPath)
	SendMessage(CHATID_LORD, "Generated story:\n"+elevenLabsPrompt)
	return nil
}

func handlerStamoce(update tgbotapi.Update) error {
	// Telegram supports only up to 10 options
	pollChoices := []string{
		"Lunneddì sera",
		"Marteddì sera",
		"Mercoleddì sera",
		"Gioveddì sera",
		"Venerdì sera",
		"Sabato mattina",
		"Sabato sera",
		"Domenica mattina",
		"Domenica sera",
		utils.PickFromArray([]string{
			"Sono calvo",
			"Oh no, sono stato Matteato",
			"Non mi sento bene, devo riposare",
			"Sesso papà Gaetano",
			"Corro nudo urlando per le strade di Napoli",
			"Grottesco",
			"Sasso",
			"C'è la partita del Napoli",
			"Simpo per Matteo Criccomoro",
			"Preparo due crostate",
			"Sono il Re dei Simp",
		}),
	}

	pollConfig := tgbotapi.SendPollConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.Message.Chat.ID,
		},
		Question:              "Quanno ce stamo?",
		Options:               pollChoices,
		IsAnonymous:           false,
		AllowsMultipleAnswers: true,
		Explanation:           "t'o devo pure spiegà?",
	}

	tgClient.Send(pollConfig)
	return nil
}

func sendAudio(chatId int64, audioPath string) {
	audioFile := tgbotapi.FilePath(audioPath)
	msg := tgbotapi.NewAudio(chatId, audioFile)

	sentMessage, err := tgClient.Send(msg)

	if err != nil && strings.Contains(err.Error(), "nil") {
		forwardMsg := tgbotapi.NewForward(CHATID_LORD, sentMessage.Chat.ID, sentMessage.MessageID)
		tgClient.Send(forwardMsg)
	} else {
		log.Printf("Couldn't forward the message with audio %s.\n%v", audioPath, err)
	}
}

func SendMessage(chatId int64, text string) {
	params := &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	}

	tgClient.SendMessage(bgCtx, params)
}
