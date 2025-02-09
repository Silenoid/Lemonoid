package piper

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
)

func GenerateVoiceNarration(prompt string) (string, error) {
	audioTitle := elevenlabs.MakeAudioTitle(prompt)
	generatedAudioFilename := audioTitle + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".mp3"
	generatedAudioDir := filepath.Join(".temp", "piper_generated")
	generatedAudioCompletePath := filepath.Join(generatedAudioDir, generatedAudioFilename)

	os.MkdirAll(generatedAudioDir, os.ModePerm)

	echoCmd := exec.Command(
		"echo",
		// strings.ReplaceAll(prompt, "\n", ""),
		"novabbe wagooooooooo sesso papa gaetano",
	)

	piperCmd := exec.Command(
		"piper-tts",
		"--model", "/home/sileno/Test/paola.onx",
		"--output-file", generatedAudioCompletePath,
		"--sentence_silence", "0.4",
	)

	piperCmd.Stdin, _ = echoCmd.StdoutPipe()

	var out bytes.Buffer
	piperCmd.Stdout = &out

	err := echoCmd.Run()
	if err != nil {
		log.Println("Error during echo command in tts generation: " + err.Error())
		return "", err
	}

	err = piperCmd.Run()
	if err != nil {
		log.Println("Error during piper command in tts generation: " + err.Error())
		return "", err
	}
	log.Println(string(out.String()))

	return generatedAudioCompletePath, nil
}
