package openai

import (
	"context"
	"log"

	"github.com/Silenoid/Lemonoid/internal/utils"
	openaiapi "github.com/sashabaranov/go-openai"
)

const OPENAI_MODEL = openaiapi.GPT4oMini

var token string
var openaiclient *openaiapi.Client

func Initialize() {
	token = utils.TokenOpenAi

	openaiclient = openaiapi.NewClient(token)
}

func GenerateStory(prompt string) (string, error) {
	if len(token) == 0 {
		panic("OpenAi token is not set")
	}

	resp, err := openaiclient.CreateChatCompletion(
		context.Background(),
		openaiapi.ChatCompletionRequest{
			Model:     OPENAI_MODEL,
			MaxTokens: 200,
			Messages: []openaiapi.ChatCompletionMessage{
				{
					Role:    openaiapi.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 1.0,
		},
	)

	if err != nil {
		log.Printf("Failing OpenAI call -> %v", err)
		return "", err
	}

	log.Printf("[OpenAi client] %s", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}
