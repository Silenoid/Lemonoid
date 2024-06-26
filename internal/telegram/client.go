package telegram

import (
	"log"
	"strings"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/history"
	"github.com/Silenoid/Lemonoid/internal/openai"
	"github.com/Silenoid/Lemonoid/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string
var tgclient *tgbotapi.BotAPI
var startTime time.Time

var CHATID_LORD int64 = 449697032
var CHATID_CHICECE int64 = -1001623264158

var PROMPT_THEMES []string = []string{
	"una avventura fantastica",
	"una lettera scritta in epoca vittoriana",
	"una sceneggiatura di uno spettacolo comico",
	"una storia di Natale",
	"un vecchio articolo di giornale",
	"una missiva in epoca medievale",
	"un antico editto romano",
	"una storia dell'orrore",
	"un racconto in un linguaggio grezzo primitivo",
	"una pagina di un diario segreto, cominciando con \"Caro diario\"",
	"un discorso politico",
	"un sermone fatto in chiesa",
	"un romanzo rosa",
	"un articolo scientifico",
	"una rivista di gossip",
	"un canto della divina commedia",
	"una poesia ermetica",
	"un libro per bambini",
}

func Initialize(isDebugging bool) {
	token = utils.TokenTelegram

	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Error during bot initialization -> ", err)
	}

	tgclient = tgbot
	log.Println("Telegram bot client authorized under the account " + tgclient.Self.UserName)

	tgclient.Debug = isDebugging

	startTime = time.Now()
	sendMessage(CHATID_LORD, "Lemonoid awakened at "+startTime.String())
}

func Listen() {
	if len(token) == 0 {
		panic("Telegram token is not set")
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updatesChannel := tgclient.GetUpdatesChan(updateConfig)
	updatesChannel.Clear()

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
	sendMessage(update.Message.Chat.ID, `Aò a manzo, eccote du seppie e ttre ppiovre de aiuto:
	/help	l'hai usato mò a cojone, ma che sei frocio?
	/tldr	azzì questo teggenera er tuloddonrì
	/status	je chiedi mammamiacommestaaa
	`)
	return nil
}

func handlerStatus(update tgbotapi.Update) error {
	ElevenLabsSubStatus := elevenlabs.GetSubscriptionStatus()
	// TODO: get OpenAI usage with a request (see openai client.go)
	sendMessage(update.Message.Chat.ID, ElevenLabsSubStatus)
	return nil
}

func handlerTldr(update tgbotapi.Update) error {
	chatHistory := history.GetChatHistory(update.Message.Chat.ID)

	openAiPromptBuilder := strings.Builder{}
	openAiPromptBuilder.WriteString("Genera un riassunto della seguente chat")
	pickedPromptTheme := utils.PickFromArray(PROMPT_THEMES)
	openAiPromptBuilder.WriteString("Genera un riassunto della seguente chat come se fosse ")
	openAiPromptBuilder.WriteString(pickedPromptTheme)
	openAiPromptBuilder.WriteString(", utilizzando almeno una volta il termine 'grottesco':\n")
	openAiPromptBuilder.WriteString(chatHistory)
	openAiPrompt := openAiPromptBuilder.String()

	log.Printf("Prompt a tema '%s': %s", pickedPromptTheme, openAiPrompt)
	elevenLabsPrompt, err := openai.GenerateStory(openAiPrompt)
	if err != nil {
		return err
	}

	generatedAudioPath, err := elevenlabs.GenerateVoiceNarration(elevenLabsPrompt)
	if err != nil {
		sendMessage(update.Message.Chat.ID, "Errore durante la generazione vocale: "+err.Error())
		return err
	}

	sendMessage(update.Message.Chat.ID, "Tema utilizzato per il prompt: "+pickedPromptTheme)
	sendAudio(update.Message.Chat.ID, generatedAudioPath)
	sendMessage(CHATID_LORD, "Generated story:\n"+elevenLabsPrompt)
	return nil
}

func handlerStamoce(update tgbotapi.Update) error {
	rightNow := time.Now()

	// Telegram supports only up to 10 options
	pollChoices := []string{
		rightNow.Weekday().String() + " sera",
		rightNow.Add(time.Hour*24*1).Weekday().String() + " mattina",
		rightNow.Add(time.Hour*24*1).Weekday().String() + " sera",
		rightNow.Add(time.Hour*24*2).Weekday().String() + " mattina",
		rightNow.Add(time.Hour*24*2).Weekday().String() + " sera",
		rightNow.Add(time.Hour*24*3).Weekday().String() + " mattina",
		rightNow.Add(time.Hour*24*3).Weekday().String() + " sera",
		rightNow.Add(time.Hour*24*4).Weekday().String() + " mattina",
		rightNow.Add(time.Hour*24*4).Weekday().String() + " sera",
		utils.PickFromArray([]string{
			"Sono calvo",
			"Oh no, sono stato Matteato",
			"Non mi sento bene, devo riposare",
			"Sesso papà Gaetano",
			"Corro nudo urlando per le strade di Napoli",
			"Grottesco",
			"Sasso",
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

	tgclient.Send(pollConfig)
	return nil
}

func sendAudio(chatId int64, audioPath string) {
	audioFile := tgbotapi.FilePath(audioPath)
	msg := tgbotapi.NewAudio(chatId, audioFile)
	tgclient.Send(msg)
}

func sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	tgclient.Send(msg)
}
