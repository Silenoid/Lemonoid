package elevenlabs

import (
	"testing"
	"time"

	"github.com/Silenoid/Lemonoid/internal/elevenlabs"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	firstTime := time.Now()
	time.Sleep(time.Second * 2)
	secondTime := time.Now()

	assert.True(t, secondTime.After(firstTime))
}

func TestClearString(t *testing.T) {
	inputStr := "Enzo: \"In poche parole, non ci capisco niente\""
	title := elevenlabs.MakeAudioTitle(inputStr)

	assert.Equal(t, "EnzoInpocheparolenon", title)
}
