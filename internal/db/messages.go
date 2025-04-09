package db

import (
	"database/sql"
)

// LogMessage inserts a new message into the database.
func LogMessage(db *sql.DB, message Message) error {
	query := `INSERT INTO messages (message_id, timestamp, name, last_name, username, group_id, user_id, content)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, message.MessageID, message.Timestamp, message.Name, message.LastName, message.Username, message.GroupID, message.UserID, message.Content)
	return err
}

// GetMessages retrieves messages from the database based on message ID and group ID.
func GetMessages(db *sql.DB, messageID int64, groupID int64) ([]Message, error) {
	query := `SELECT message_id, timestamp, name, last_name, username, group_id, user_id
              FROM messages
              WHERE group_id = ? AND message_id >= ?
              ORDER BY timestamp ASC LIMIT 500`

	rows, err := db.Query(query, groupID, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.MessageID, &msg.Timestamp, &msg.Name, &msg.LastName, &msg.Username, &msg.GroupID, &msg.UserID, &msg.Content); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
