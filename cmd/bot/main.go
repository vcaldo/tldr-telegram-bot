package main

import (
	"log"

	"tldr-telegram-bot/internal/config"
	"tldr-telegram-bot/internal/db"
	"tldr-telegram-bot/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation error: %v", err)
	}

	// Initialize database
	db.InitDB()

	// Start the Telegram bot
	bot, err := telegram.NewBot()
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	log.Println("Bot started and listening for messages...")
	bot.Start()
}


cron.AddFunc("0 3 * * *", func() {
	telegram.RunDailySummary()
})

