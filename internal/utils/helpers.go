package utils

import (
    "strings"
)

// TrimAndLower trims whitespace from the input string and converts it to lowercase.
func TrimAndLower(input string) string {
    return strings.ToLower(strings.TrimSpace(input))
}

// IsEmpty checks if the input string is empty.
func IsEmpty(input string) bool {
    return strings.TrimSpace(input) == ""
}

// Contains checks if a slice of strings contains a specific string.
func Contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// JoinMessages formats a slice of messages into a single string with user identifiers.
func JoinMessages(messages []string, userIdentifier string) string {
    return strings.Join(messages, "\n"+userIdentifier+": ")
}