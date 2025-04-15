package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"

	"tldr-telegram-bot/internal/config"
	"tldr-telegram-bot/internal/db"
	"tldr-telegram-bot/internal/llm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var triggerWords = []string{"resuma", "resume", "tldr", "summary", "toguro por favor", "toguro please", "professor toguro", "professor toguro por favor", "professor toguro please", "toguro", "toguro por favor", "toguro please", "toguro professor", "toguro professor por favor", "toguro professor please"}

func HandleMessage(update tgbotapi.Update) {
	if update.Message == nil || update.Message.ReplyToMessage == nil {
		return
	}

	if !isAuthorizedGroup(update.Message.Chat.ID) {
		logUnauthorizedAttempt(update.Message.Chat.ID)
		return
	}

	if isTriggerWord(update.Message.Text) {
		log.Printf("Trigger word detected in group %d", update.Message.Chat.ID)
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
	concatenatedText = strings.ReplaceAll(concatenatedText, "\n", " ")
	concatenatedText = strings.TrimSpace(concatenatedText)
	fmt.Println("Concatenated text for summarization:", concatenatedText)

	if os.Getenv("LOCAL_MODEL") != "true" {
		summary, err := llm.SummarizeGemini(concatenatedText, myConfig.Lang)
		if err != nil {
			log.Printf("Error summarizing messages: %v", err)
			return
		}
		sendSummary(update.Message.Chat.ID, summary)
		return
	}
	summary, err := llm.Summarize(concatenatedText, myConfig.Lang)
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
		switch {
		case msg.Name != "" && msg.LastName != "":
			sender = fmt.Sprintf("%s %s", msg.Name, msg.LastName)
		case msg.Name != "":
			sender = msg.Name
		case msg.LastName != "":
			sender = msg.LastName
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", sender, msg.Content))
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

func summarizeDay(chatID int64) {
	myDb := db.GetDB()
	myConfig, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Printf("Error loading time location: %v", err)
		loc = time.FixedZone("BRT", -3*60*60) // fallback
	}

	now := time.Now().In(loc)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, loc)

	var finalSummary strings.Builder

	for i := 0; i < 8; i++ {
		from := dayStart.Add(time.Duration(i*3) * time.Hour)
		to := from.Add(3 * time.Hour)

		messages, err := db.GetMessagesByTimeRange(myDb, chatID, from, to)
		if err != nil {
			log.Printf("Error fetching messages for range %s - %s: %v", from.Format("15:04"), to.Format("15:04"), err)
			continue
		}

		if len(messages) == 0 {
			continue // skip empty blocks
		}

		text := formatMessages(messages)
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.TrimSpace(text)

		var summary string
		if os.Getenv("LOCAL_MODEL") != "true" {
			summary, err = llm.SummarizeGemini(text, myConfig.Lang)
		} else {
			summary, err = llm.Summarize(text, myConfig.Lang)
		}

		if err != nil {
			log.Printf("Error summarizing block %s - %s: %v", from.Format("15:04"), to.Format("15:04"), err)
			continue
		}

		finalSummary.WriteString(fmt.Sprintf("🕒 %s - %s\n%s\n\n", from.Format("15:04"), to.Format("15:04"), summary))
	}

	if finalSummary.Len() > 0 {
		sendSummary(chatID, finalSummary.String())
	} else {
		log.Printf("No content to summarize for group %d", chatID)
	}
}



func RunDailySummary() {
	myConfig, err := config.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	for _, groupID := range myConfig.AuthorizedGroups {
		log.Printf("Running daily summary for group: %d", groupID)
		summarizeDay(groupID)
	}
}