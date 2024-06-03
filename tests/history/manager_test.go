package history

import (
	"testing"

	"github.com/Silenoid/Lemonoid/internal/history"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	history.CleanChatMap()
}

func TestGenerateSummary_BasicGeneration(t *testing.T) {
	history.AddMessageToChatHistory(1, 1, "pippo", "no")
	history.AddMessageToChatHistory(1, 1, "pippo", "va")
	history.AddMessageToChatHistory(1, 1, "pippo", "be")
	expectedSummary := "pippo: \"no\".\npippo: \"va\".\npippo: \"be\".\n"

	actualSummary := history.GetChatHistory(1)

	assert.Equal(t, expectedSummary, actualSummary)
	t.Cleanup(cleanup)
}

func TestGenerateSummary_UseRealUsernames(t *testing.T) {
	history.AddMessageToChatHistory(1, 449697032, "", "no")
	history.AddMessageToChatHistory(1, 76783158, "", "va")
	history.AddMessageToChatHistory(1, 595335090, "", "be")
	expectedSummary := "Enzo: \"no\".\nMatteo finto: \"va\".\nGraizano: \"be\".\n"

	actualSummary := history.GetChatHistory(1)

	assert.Equal(t, expectedSummary, actualSummary)
	t.Cleanup(cleanup)
}

func TestGenerateSummary_UnknownName(t *testing.T) {
	history.AddMessageToChatHistory(1, 123123, "", "no")
	history.AddMessageToChatHistory(1, 234234, "", "va")
	history.AddMessageToChatHistory(1, 345345, "", "be")
	expectedSummary := "123123: \"no\".\n234234: \"va\".\n345345: \"be\".\n"

	actualSummary := history.GetChatHistory(1)

	assert.Equal(t, expectedSummary, actualSummary)
	t.Cleanup(cleanup)
}
