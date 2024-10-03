package utils

import (
	"testing"
	"time"

	"github.com/Silenoid/Lemonoid/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestToReadableDate_BasicGeneration(t *testing.T) {
	actual := utils.ToReadableDate(time.Now())
	assert.Regexp(t, "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}", actual)
}

func TestToReadableSince_BasicGeneration(t *testing.T) {
	actual := utils.ToReadableSince(time.Now(), time.Now().Add(-20*time.Hour).Add(-33*time.Minute))
	assert.Equal(t, "20h33m0s", actual)
}

func TestToReadableHowLongTo_BasicGeneration(t *testing.T) {
	actual := utils.ToReadableHowLongTo(time.Now(), time.Now().Add(-1*time.Hour), 2*time.Hour)
	assert.Equal(t, "1h0m0s", actual)
}
