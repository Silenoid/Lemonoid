package telegram

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/gemini"
	"github.com/Silenoid/Lemonoid/internal/history"
	"github.com/Silenoid/Lemonoid/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const TLDRPRO_WAIT_TIME = time.Hour * 72

var lastTLDRPROTime time.Time
var isFirstTLDRPRO = false

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		utils.PrintMsg(update)

		if time.Now().Before(startTime.Add(3 * time.Second)) {
			log.Printf("Not processing message [%b - %s] because of previous messages cleanup",
				update.Message.Chat.ID,
				update.Message.From.Username,
			)
			return
		}

		if isAllowedChatId(update.Message.Chat.ID) {
			if update.Message.Text != "" {
				log.Printf(
					"Remembering message from [%b - %s] with ID %b",
					update.Message.Chat.ID,
					update.Message.From.Username,
					update.Message.ID,
				)
				history.AddMessageToChatHistory(
					update.Message.Chat.ID,
					update.Message.From.ID,
					update.Message.From.Username,
					update.Message.Text,
				)
			}
		} else {
			log.Printf(
				"Skkipping message from [%b - %s] with ID %b",
				update.Message.Chat.ID,
				update.Message.From.Username,
				update.Message.ID,
			)
		}
	}
}

func handlerHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	RespondWithText(update, `Aò a manzo, eccote du seppie e ttre ppiovre de aiuto:
	/help	l'hai usato mò a cojone, ma che sei frocio?
	/tldr	azzì questo teggenera er tuloddonrì
	/status	je chiedi mammamiacommestaaa
	`)
}

func handlerStatus(ctx context.Context, b *bot.Bot, update *models.Update) {
	ElevenLabsSubStatus := elevenlabs.GetSubscriptionStatus()
	RespondWithText(update, ElevenLabsSubStatus)
}

func handlerTldrPro(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !isFirstTLDRPRO && time.Now().Before(lastTLDRPROTime.Add(TLDRPRO_WAIT_TIME)) {
		RespondWithText(update, "A cojò, li mortacci stracci tua ma che me voj rovinà? So ppasati solo "+utils.ToReadableSince(time.Now(), lastTLDRPROTime)+" da nantro vocale. Statte bbono pe' n'antri "+utils.ToReadableHowLongTo(time.Now(), lastTLDRPROTime, TLDRPRO_WAIT_TIME))
		return
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
	RespondWithText(update, "Tema utilizzato per il prompt: "+pickedPromptTheme)

	generatedStory, err := gemini.GenerateStory(llmPrompt)
	if err != nil {
		if strings.Contains(update.Message.Text, "exceeded your current quota") {
			RespondWithText(update, "Ao, so ffiniti li sordi pe generà er testo li mortacci stracci")
		}
		SendMessage(CHATID_LORD, err.Error())
	}

	lastTLDRPROTime = time.Now()
	isFirstTLDRPRO = false

	generatedAudioPath, err := elevenlabs.GenerateVoiceNarration(generatedStory, pickedVoice)
	if err != nil {
		RespondWithText(update, "Errore nella generazione vocale, dunque beccate solo er testo generato e muto:\n"+generatedStory)
	}

	sendAudio(update, generatedAudioPath)
	SendMessage(CHATID_LORD, "Generated story:\n"+generatedStory)
}

func handlerTldr(ctx context.Context, b *bot.Bot, update *models.Update) {
	RespondWithText(update, "Stamo lavorannoce")
	params := &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileString{Data: "AgACAgIAAxkDAAIBOWJimnCJHQJiJ4P3aasQCPNyo6mlAALDuzEbcD0YSxzjB-vmkZ6BAQADAgADbQADJAQ"},
		Caption: "Preloaded Facebook logo",
	}

	b.SendPhoto(ctx, params)
}

func handlerStamoce(ctx context.Context, b *bot.Bot, update *models.Update) {
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
			MessageThreadID:       update.Message.MessageThreadID,
			Question:              "Quanno ce stamo?",
			Options:               pollOptions,
			IsAnonymous:           bot.False(),
			AllowsMultipleAnswers: *bot.True(),
			Explanation:           "t'o devo pure spiegà?",
		},
	)
}
