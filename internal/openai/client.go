package openai

import (
	"context"
	"errors"
	"log"

	"github.com/Silenoid/Lemonoid/internal/utils"
	openaiapi "github.com/sashabaranov/go-openai"
)

const OPENAI_MODEL = openaiapi.GPT3Dot5Turbo

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
			Model: OPENAI_MODEL,
			Messages: []openaiapi.ChatCompletionMessage{
				{
					Role:    openaiapi.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 1.5,
		},
	)

	if err != nil {
		log.Printf("Failing OpenAI call -> %v", err)
		return "", errors.New("error during OpenAi call")
	}

	log.Printf("[OpenAi client] %s", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content, nil
}

// TODO: retrieve api usage
// https://stackoverflow.com/questions/75703189/api-to-get-current-balance-tokens
// https://api.openai.com/dashboard/billing/credit_grants
// https://pkg.go.dev/net/http
