package db

import "time"

// Message represents the schema for storing messages in the database.
type Message struct {
	MessageID int64     `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
}
