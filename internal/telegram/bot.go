package telegram

import (
	"log"
	"tldr-telegram-bot/internal/db"

	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot() (*Bot, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
		return nil, err
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
		return nil, err
	}

	api.Debug = false
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{api: api}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-message updates
			continue
		}

		// Log incoming messages
		go b.logMessage(update.Message)

		// handle incoming messages
		go HandleMessage(update)
	}
}

func (b *Bot) logMessage(message *tgbotapi.Message) {
	parsedMsg := db.Message{
		MessageID: int64(message.MessageID),
		Timestamp: message.Time(),
		Name:      message.From.FirstName,
		LastName:  message.From.LastName,
		Username:  message.From.UserName,
		GroupID:   message.Chat.ID,
		UserID:    message.From.ID,
		Content:   message.Text,
	}

	db.InitDB()
	myDb := db.GetDB()
	if myDb == nil {
		log.Println("Database connection is nil")
		return
	}

	if err := db.LogMessage(myDb, parsedMsg); err != nil {
		log.Printf("failed to insert message: %v", err)
	}
}
