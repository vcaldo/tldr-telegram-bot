package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Validate checks the required environment variables and their values.
func Validate() error {
	requiredVars := []string{
		"TELEGRAM_BOT_TOKEN",
		"DEFAULT_LANG",
		"OLLAMA_MODEL",
		"AUTHORIZED_GROUPS",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return errors.New("missing required environment variable: " + v)
		}
	}

	// Validate DEFAULT_LANG
	defaultLang := os.Getenv("DEFAULT_LANG")
	if !isValidLanguage(defaultLang) {
		return errors.New("invalid DEFAULT_LANG value: " + defaultLang)
	}

	// Validate AUTHORIZED_GROUPS
	if err := validateAuthorizedGroups(os.Getenv("AUTHORIZED_GROUPS")); err != nil {
		return err
	}

	return nil
}

// isValidLanguage checks if the provided language is one of the allowed values.
func isValidLanguage(lang string) bool {
	allowedLanguages := []string{"pt", "en", "es"}
	for _, l := range allowedLanguages {
		if lang == l {
			return true
		}
	}
	return false
}

// validateAuthorizedGroups checks if the authorized groups are valid numeric IDs.
func validateAuthorizedGroups(groups string) error {
	groupIDs := strings.Split(groups, ",")
	for _, id := range groupIDs {
		if _, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64); err != nil {
			return errors.New("invalid group ID in AUTHORIZED_GROUPS: " + id)
		}
	}
	return nil
}