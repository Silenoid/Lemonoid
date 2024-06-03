package history

import (
	"log"
	"strconv"
	"strings"
)

var chatMap map[int64][]string = make(map[int64][]string)

var trueUsernames = map[int64]string{
	449697032: "Enzo",
	76783158:  "Matteo finto",
	595335090: "Graizano",
	344356561: "Alessandro",
	380320701: "Novy",
	179330209: "Simone",
	865721543: "Roberto",
	466018597: "Paolo",
}

const CHAT_MSG_THRESHOLD = 150

func AddMessageToChatHistory(chatId int64, senderId int64, senderName string, messageText string) {
	chatEntryText := getUserName(senderId, senderName) + ": \"" + messageText + "\".\n"
	if chatEntries, exists := chatMap[chatId]; exists {
		chatMap[chatId] = append(chatEntries, chatEntryText)
	} else {
		chatMap[chatId] = []string{chatEntryText}
	}
}

func getUserName(senderId int64, senderName string) string {
	if trueUsername, exists := trueUsernames[senderId]; exists {
		return trueUsername
	} else if senderName != "" {
		return senderName
	} else {
		log.Printf("Warning: username %d not recognized. Using just its ID", senderId)
		return strconv.Itoa(int(senderId))
	}
}

func GetChatHistory(chatId int64) string {
	summary := strings.Builder{}

	chatEntries := chatMap[chatId]
	oldestEntryIndex := max(0, len(chatEntries)-CHAT_MSG_THRESHOLD)
	latestEntryIndex := len(chatEntries)

	for i := oldestEntryIndex; i < latestEntryIndex; i++ {
		summary.WriteString(chatEntries[i])
	}

	return summary.String()
}

func CleanChatMap() {
	chatMap = map[int64][]string{}
}
