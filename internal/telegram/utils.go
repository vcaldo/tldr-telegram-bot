package telegram

import (
	"fmt"
	"strings"
)

// FormatMessage formats a message with the user's identifier.
func FormatMessage(name, lastName, username string, userID int64) string {
	if name != "" && lastName != "" {
		return fmt.Sprintf("%s %s: ", name, lastName)
	} else if name != "" {
		return fmt.Sprintf("%s: ", name)
	} else if username != "" {
		return fmt.Sprintf("@%s: ", username)
	}
	return fmt.Sprintf("%d: ", userID)
}

// IsTriggerWord checks if the message contains any of the trigger words.
func IsTriggerWord(message string) bool {
	triggerWords := []string{"resuma", "tldr", "summary", "toguro por favor", "toguro please", "toguro"}
	for _, word := range triggerWords {
		if strings.EqualFold(strings.TrimSpace(message), word) {
			return true
		}
	}
	return false
}