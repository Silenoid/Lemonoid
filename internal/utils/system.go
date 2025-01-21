package utils

import (
	"math/rand"
	"os"
)

var TokenTelegram string
var TokenOpenAi string
var TokenElevenLabs string
var TokenDiscord string

func LoadEnvVars() {
	TokenTelegram = GetAndValidateEnvVar("LEMONOID_TOKEN_TELEGRAM")
	TokenOpenAi = GetAndValidateEnvVar("LEMONOID_TOKEN_OPENAI")
	TokenElevenLabs = GetAndValidateEnvVar("LEMONOID_TOKEN_ELEVENLABS")
	TokenDiscord = GetAndValidateEnvVar("LEMONOID_TOKEN_DISCORD")
}

func GetAndValidateEnvVar(envVarName string) string {
	value, exists := os.LookupEnv(envVarName)
	if exists {
		return value
	} else {
		panic("Exiting due to missing environment variable " + envVarName)
	}
}

func PickFromArray(arrayToPickFrom []string) string {
	return arrayToPickFrom[rand.Intn(len(arrayToPickFrom))]
}
