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

	piperCmd := exec.Command(
		"piper-tts",
		"--model", "/home/sileno/Test/paola.onnx",
		"--output-file", generatedAudioCompletePath,
		"--sentence_silence", "0.4",
	)

	piperCmd.Stdin = strings.NewReader(prompt)

	eout, err := piperCmd.CombinedOutput()
	if err != nil {
		log.Println("Error during piper command in tts generation: " + err.Error())
		return "", err
	}
	log.Println(eout[:])

	return generatedAudioCompletePath, nil
}
