package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SaveMessage(message *models.Message) error {
	query := `INSERT INTO messages (sender, receiver, content, timestamp) VALUES (?, ?, ?, ?)`
	_, err := repo.db.Exec(query, message.Sender, message.Receiver, message.Content, message.Timestamp)
	return err
}

func (repo *MessageRepository) GetMessages(sender, receiver string) ([]models.Message, error) {
	var messages []models.Message
	query := `SELECT * FROM messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) ORDER BY timestamp`
	err := repo.db.Select(&messages, query, sender, receiver, receiver, sender)
	return messages, err
}

// GetMessageByID retrieves a message by ID from the database
func (repo *MessageRepository) GetMessageByID(id string) (models.Message, error) {
	var message models.Message
	row := repo.db.QueryRow(`SELECT id, sender, receiver, content, timestamp, is_forwarded, original_sender, original_message_id
		FROM messages WHERE id = ?`, id)

	err := row.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.Timestamp, &message.IsForwarded, &message.OriginalSender, &message.OriginalMessageID)
	if err == sql.ErrNoRows {
		return message, nil // No error, but no record found
	} else if err != nil {
		return message, err
	}
	return message, nil
}
