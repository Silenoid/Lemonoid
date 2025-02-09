package piper

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
		strings.ReplaceAll(prompt, "\n", ""),
	)

	piperCmd := exec.Command(
		"piper-tts",
		"--model", "/home/sileno/Test/paola.onx",
		"--output-file", generatedAudioCompletePath,
		"--sentence_silence", "0.4",
	)

	echoPipe, _ := echoCmd.StdoutPipe()
	defer echoPipe.Close()

	piperCmd.Stdin = echoPipe

	echoCmd.Start()
	stdout, err := piperCmd.CombinedOutput()
	if err != nil {
		log.Println("Error during piper tts generation: " + err.Error())
		return "", err
	}
	log.Println(string(stdout[:]))

	return generatedAudioCompletePath, nil
}
