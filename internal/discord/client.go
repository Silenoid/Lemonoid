package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Silenoid/Lemonoid/internal/telegram"
	"github.com/Silenoid/Lemonoid/internal/utils"
	"github.com/bwmarrin/discordgo"
)

const GUILD_ID = "830811876761010196"

func Initialize() {
	var token = utils.TokenElevenLabs
	discordClient, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("[Discord client] Error creating Discord session: %s", err)
		return
	}

	// Called every time a new message is created on any guild that the authenticated bot has access to
	discordClient.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log.Printf("[Discord client] %s says: %s", m.Author.ID, m.Content)

		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		} else if m.Content == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}
	})

	// Called every time an event is related to a voice chat
	discordClient.AddHandler(func(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
		guild, err := s.State.Guild(GUILD_ID)
		if err != nil {
			log.Printf("[Discord client] Error fetching guild: %s", err)
			return
		}

		if len(guild.VoiceStates) > 0 {
			telegram.SendMessage(telegram.CHATID_LORD, "ao anvedi")
		}

	})

	// We only care about receiving message events.
	discordClient.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open a websocket connection to Discord and begin listening.
	err = discordClient.Open()
	if err != nil {
		fmt.Println("Discord client: error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discordClient.Close()
}
