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

	cmd := exec.Command(
		"echo",
		"'"+strings.ReplaceAll(prompt, "\n", "")+"'",
		" | ",
		"piper-tts",
		"--model", "/home/sileno/Test/paola.onx",
		"--output-file", generatedAudioCompletePath,
	)

	stdout, err := cmd.Output()
	if err != nil {
		log.Println("Error during piper tts generation: " + err.Error())
		return "", err
	}
	log.Println(string(stdout[:]))

	return generatedAudioCompletePath, nil
}
