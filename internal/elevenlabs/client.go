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

var CLONED_VOICES []string = []string{
	"i86lB8eIKMQcO470EIFz", // G
	"ml5JfpB48j688Rpbbz2M", // G Maronne
	"d9Gr3L3YR4d9Sf9Gt8cV", // S
	"IzoLtTXseyrunESwWmw3"} // M

var BASIC_VOICES []string = []string{
	"EXAVITQu4vr4xnSDxMaL",
	"FGY2WhTYpPnrIDTdsKH5",
	"IKne3meq5aSn9XLyUdCD",
	"JBFqnCBsd6RMkjVDRZzb",
	"N2lVS1w4EtoT3dr4eOWO",
	"TX3LPaxmHKxFdv7VOQHJ",
	"XB0fDUnXU5powFXDhCwa",
	"Xb7hH8MSUJpSbSDYk0k2",
	"XrExE9yKIg1WjnnlVkGX",
	"bIHbv24MWmeRgasZH58o",
	"cgSgspJ2msm6clMCkdW9",
	"cjVigY5qzO86Huf0OWal",
	"iP95p4xoKVk53GoZ742B",
	"nPczCjzI2devNBz1zQrb",
	"onwK4e9ZLuTAKqWW03F9",
	"pFZP5JQG7iQjIQuC4Bku",
	"pqHfZKP75CvOlQylNhV4"}

const GENERATION_WAITING_PERIOD = time.Hour * 24

var token string
var elclient *elapi.Client
var lastGeneratedAudioTime time.Time
var isFirstGeneration = false

func Initialize() {
	token = utils.TokenElevenLabs

	elclient = elapi.NewClient(context.Background(), token, 30*time.Second)
}

func GenerateVoiceNarration(prompt string, pickedVoice string) (string, error) {
	if len(token) == 0 {
		panic("ElevenLabs token is not set")
	}

	if isFirstGeneration || time.Now().After(lastGeneratedAudioTime.Add(GENERATION_WAITING_PERIOD)) {
		payload := elapi.TextToSpeechRequest{
			Text:    prompt,
			ModelID: "eleven_multilingual_v2",
		}

		audio, err := elclient.TextToSpeech(pickedVoice, payload)
		if err != nil {
			log.Printf("Failing ElevenLabs call -> %v", err)
		}

		lastGeneratedAudioTime = time.Now()
		isFirstGeneration = false

		audioTitle := MakeAudioTitle(prompt)
		generatedAudioFilename := audioTitle + "-" + strconv.FormatInt(lastGeneratedAudioTime.UnixMilli(), 10) + ".mp3"
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
		return "", errors.New("A cojò, li mortacci stracci tua ma che me voj rovinà? So ppasati solo " + utils.ToReadableSince(time.Now(), lastGeneratedAudioTime) + " da nantro vocale. Statte bbono pe' n'antri " + utils.ToReadableHowLongTo(time.Now(), lastGeneratedAudioTime, GENERATION_WAITING_PERIOD))
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
