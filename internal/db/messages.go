package db

import (
	"database/sql"
	"log"
	"time"
)

// LogMessage inserts a new message into the database.
func LogMessage(db *sql.DB, message Message) error {
	query := `INSERT INTO messages (message_id, timestamp, name, last_name, username, group_id, user_id, content)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query,
		message.MessageID,
		message.Timestamp,
		message.Name,
		message.LastName,
		message.Username,
		message.GroupID,
		message.UserID,
		message.Content,
	)

	return err
}

// GetMessages retrieves messages from the database based on message ID and group ID.
func GetMessages(db *sql.DB, messageID int64, groupID int64) ([]Message, error) {
	firstMessageTimestamp, err := getMessageTimestamp(db, messageID, groupID)
	if err != nil {
		return nil, err
	}

	query := `SELECT message_id, timestamp, name, last_name, username, group_id, user_id, content
		  FROM messages
		  WHERE group_id = $1 AND timestamp BETWEEN $2 AND ($2 + interval '10 minutes')
		  ORDER BY timestamp ASC LIMIT 1000`

	log.Default().Printf("Query: %s, GroupID: %d, Timestamp: %v", query, groupID, firstMessageTimestamp)
	rows, err := db.Query(query, groupID, firstMessageTimestamp)
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

func getMessageTimestamp(db *sql.DB, messageID int64, groupID int64) (*time.Time, error) {
	query := `SELECT timestamp FROM messages WHERE message_id = $1 AND group_id = $2`
	row := db.QueryRow(query, messageID, groupID)
	var timestamp time.Time
	if err := row.Scan(&timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No message found
		}
		return nil, err
	}
	return &timestamp, nil
}
