package telegram

import (
	"log"
	"strings"

	"tldr-telegram-bot/internal/config"
	"tldr-telegram-bot/internal/db"
	"tldr-telegram-bot/internal/llm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MAX_MESSAGES = 500

var triggerWords = []string{"resuma", "tldr", "summary", "toguro por favor", "toguro please", "toguro"}

func HandleMessage(update tgbotapi.Update) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return
	}

	if !isAuthorizedGroup(update.Message.Chat.ID) {
		logUnauthorizedAttempt(update.Message.Chat.ID)
		return
	}

	if isTriggerWord(update.Message.ReplyToMessage.Text) {
		collectAndSummarizeMessages(update)
	}
}

func isAuthorizedGroup(groupID int64) bool {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return false
	}

	// Check if the group ID is in the list of authorized groups
	for _, authorizedGroup := range config.AuthorizedGroups {
		if authorizedGroup == groupID {
			return true
		}
	}
	log.Printf("Unauthorized group access attempt: %d", groupID)
	return false
}

func logUnauthorizedAttempt(groupID int64) {
	log.Printf("Unauthorized access attempt in group: %d", groupID)
}

func isTriggerWord(message string) bool {
	for _, word := range triggerWords {
		if strings.Contains(strings.ToLower(message), word) {
			return true
		}
	}
	return false
}

func collectAndSummarizeMessages(update tgbotapi.Update) {
	myDb := db.GetDB()

	myConfig, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	messages, err := db.GetMessages(myDb, int64(update.Message.ReplyToMessage.MessageID), update.Message.Chat.ID)
	if err != nil {
		log.Printf("Error collecting messages: %v", err)
		return
	}

	if len(messages) == 0 {
		log.Println("No messages found for summarization.")
		return
	}

	concatenatedText := formatMessages(messages)

	summary, err := llm.Summarize(concatenatedText, myConfig.DefaultLang)
	if err != nil {
		log.Printf("Error summarizing messages: %v", err)
		return
	}

	sendSummary(update.Message.Chat.ID, summary)
}

func formatMessages(messages []db.Message) string {
	var sb strings.Builder
	for _, msg := range messages {
		sender := msg.Username
		if sender == "" {
			if msg.Name != "" {
				sender = msg.Name
				if msg.LastName != "" {
					sender += " " + msg.LastName
				}
			} else {
				sender = string(msg.UserID)
			}
		}
		sb.WriteString(sender + ": " + msg.Content + "\n")
	}
	return sb.String()
}

func sendSummary(chatID int64, summary string) {
	bot, err := NewBot()
	if err != nil {
		log.Printf("Error creating bot instance: %v", err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, summary)
	_, err = bot.api.Send(msg)
	if err != nil {
		log.Printf("Error sending summary: %v", err)
	}
}
