package discord

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Silenoid/Lemonoid/internal/telegram"
	"github.com/Silenoid/Lemonoid/internal/utils"
	"github.com/bwmarrin/discordgo"
)

const GUILD_ID = "830811876761010196"

var discordSession *discordgo.Session
var isDeadChat = true
var timelastSentEnterMessage = time.Now().Add(-30 * time.Second)
var timelastSentDeadChat = time.Now().Add(-30 * time.Second)

func Initialize() {
	var token = utils.TokenDiscord
	newDiscordSession, err := discordgo.New("Bot " + token)
	discordSession = newDiscordSession
	if err != nil {
		log.Printf("[Discord client] Error creating Discord session: %s", err)
		return
	}

	// Called every time a new message is created on any guild that the authenticated bot has access to
	discordSession.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	discordSession.AddHandler(func(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
		guild, err := s.State.Guild(GUILD_ID)
		if err != nil {
			log.Printf("[Discord client] Error fetching guild: %s", err)
			return
		}

		if len(guild.VoiceStates) > 0 && isDeadChat {
			isDeadChat = false
			if time.Now().After(timelastSentEnterMessage.Add(30 * time.Second)) {
				telegram.SendMessage(telegram.CHATID_CHICECE, "Aò, su Discorde è cominciata a' mattanza da "+event.Member.User.Username)
				timelastSentEnterMessage = time.Now()
			}
		} else if len(guild.VoiceStates) == 0 && !isDeadChat {
			isDeadChat = true
			if time.Now().After(timelastSentDeadChat.Add(30 * time.Second)) {
				telegram.SendMessage(telegram.CHATID_CHICECE, "Te dico fermate. Su Discorde è dead chat")
				timelastSentDeadChat = time.Now()
			}
		}
	})

	// We only care about receiving message events.
	discordSession.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
}

func Listen() {
	// Open a websocket connection to Discord and begin listening.
	err := discordSession.Open()
	if err != nil {
		log.Printf("[Discord client] Error opening connection: %s", err)
		return
	} else {
		log.Println("[Discord client] Discord client is listening...")
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discordSession.Close()
}
