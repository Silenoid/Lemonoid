package gemini

import (
	"context"
	"errors"

	"github.com/Silenoid/Lemonoid/internal/utils"
	"google.golang.org/genai"
)

const GEMINI_MODEL = "gemini-2.0-flash"

var backgroundContext context.Context
var geminiClient *genai.Client

// check usage at https://aistudio.google.com/prompts/new_chat
func Initialize() {
	token := utils.TokenGemini
	backgroundContext := context.Background()

	client, err := genai.NewClient(
		backgroundContext,
		&genai.ClientConfig{
			APIKey:  token,
			Backend: genai.BackendGeminiAPI,
		},
	)

	if err != nil {
		panic(err)
	}

	geminiClient = client
}

func GenerateStory(prompt string) (string, error) {
	temperature := 1.0

	result, err := geminiClient.Models.GenerateContent(
		backgroundContext,
		GEMINI_MODEL,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature: &temperature,
		},
	)

	if err != nil {
		return "", errors.New("Error during Gemini call: " + err.Error())
	}

	return result.Text()
}
