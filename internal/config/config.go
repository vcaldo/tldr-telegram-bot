package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	Lang             string
	OllamaModel      string
	AuthorizedGroups []int64
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	authorizedGroups := os.Getenv("AUTHORIZED_GROUPS")
	groupIDs := parseAuthorizedGroups(authorizedGroups)

	return &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Lang:             os.Getenv("DEFAULT_LANG"),
		OllamaModel:      os.Getenv("OLLAMA_MODEL"),
		AuthorizedGroups: groupIDs,
	}, nil
}

func parseAuthorizedGroups(groups string) []int64 {
	var groupIDs []int64
	for _, group := range strings.Split(groups, ",") {
		groupID := strings.TrimSpace(group)
		if groupID != "" {
			id, err := strconv.ParseInt(groupID, 10, 64)
			if err != nil {
				log.Printf("Error parsing group ID %s: %v", groupID, err)
				continue
			}
			groupIDs = append(groupIDs, id)
		}
	}
	return groupIDs
}
