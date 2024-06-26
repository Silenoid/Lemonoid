package main

import (
	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/openai"
	"github.com/Silenoid/Lemonoid/internal/telegram"
	"github.com/Silenoid/Lemonoid/internal/utils"
)

func main() {
	utils.PrintWelcome()
	utils.LoadEnvVars()

	telegram.Initialize(false)
	openai.Initialize()
	elevenlabs.Initialize()

	telegram.Listen()
}
