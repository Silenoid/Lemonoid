package elevenlabs

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/Silenoid/Lemonoid/internal/utils"
	elapi "github.com/haguro/elevenlabs-go"
)

var VOICES []string = []string{
	"i86lB8eIKMQcO470EIFz", // G
	"d9Gr3L3YR4d9Sf9Gt8cV"} // S
const GENERATION_WAITING_PERIOD = time.Hour * 24

var token string
var elclient *elapi.Client
var lastGeneratedAudio time.Time
var isFirstGeneration = false

func Initialize() {
	token = utils.TokenElevenLabs

	elclient = elapi.NewClient(context.Background(), token, 30*time.Second)
}

func GenerateVoiceNarration(prompt string) (string, error) {
	if len(token) == 0 {
		panic("ElevenLabs token is not set")
	}

	if isFirstGeneration || time.Now().After(lastGeneratedAudio.Add(GENERATION_WAITING_PERIOD)) {
		payload := elapi.TextToSpeechRequest{
			Text:    prompt,
			ModelID: "eleven_multilingual_v2",
		}

		pickedVoice := utils.PickFromArray(VOICES)
		audio, err := elclient.TextToSpeech(pickedVoice, payload)
		if err != nil {
			log.Printf("Failing ElevenLabs call -> %v", err)
		}

		lastGeneratedAudio = time.Now()
		isFirstGeneration = false

		audioTitle := MakeAudioTitle(prompt)
		generatedAudioFilename := audioTitle + "-" + strconv.FormatInt(lastGeneratedAudio.UnixMilli(), 10) + ".mp3"
		generatedAudioDir := filepath.Join(".temp", "elevenlabs_generated")
		generatedAudioCompletePath := filepath.Join(generatedAudioDir, generatedAudioFilename)

		os.MkdirAll(generatedAudioDir, os.ModePerm)

		log.Printf("[ElevenLabs client] Saving elevenlabs generated audio '%s'", generatedAudioCompletePath)
		if err := os.WriteFile(generatedAudioCompletePath, audio, os.ModePerm); err != nil {
			log.Printf("Failing save ElevenLabs mp3 file writing -> %v", err)
			return "", err
		}

		return generatedAudioCompletePath, nil
	} else {
		return "", errors.New("not enough time (" + time.Since(lastGeneratedAudio).String() + ") has passed since the last generated audio. Wait for " + GENERATION_WAITING_PERIOD.String())
	}
}

func GetSubscriptionStatus() string {
	sub, err := elclient.GetSubscription()
	if err != nil {
		log.Printf("[ElevenLabs client] Problem during subscription request\n%v", err)
	}

	return "ElevenLabs: Usage at " + strconv.FormatFloat((float64(sub.CharacterCount)/float64(sub.CharacterLimit))*100.0, 'f', 2, 64) + "% (" + strconv.Itoa(sub.CharacterCount) + " of " + strconv.Itoa(sub.CharacterLimit) + ")"
}

func MakeAudioTitle(inputStr string) string {
	noSpecialCharsInput := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(inputStr, "")
	maxCharIdx := min(20, len(noSpecialCharsInput))
	return noSpecialCharsInput[0:maxCharIdx]
}
