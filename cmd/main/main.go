package main

import (
	"github.com/Silenoid/Lemonoid/internal/discord"
	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/Silenoid/Lemonoid/internal/gemini"
	"github.com/Silenoid/Lemonoid/internal/openai"
	"github.com/Silenoid/Lemonoid/internal/telegram"
	"github.com/Silenoid/Lemonoid/internal/utils"
)

func main() {
	utils.PrintWelcome()
	utils.LoadEnvVars()

	telegram.Initialize()
	openai.Initialize()
	gemini.Initialize()
	elevenlabs.Initialize()
	discord.Initialize()

	go discord.Listen()
	telegram.Listen()
}
