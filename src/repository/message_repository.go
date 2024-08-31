package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/utils"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SaveMessage(message *models.Message) error {

	var lastInsertID int64 = 0
	selectQuery := "SELECT sender, receiver , chat_id FROM messages WHERE (sender = ? and receiver = ?) OR (receiver = ? and sender = ?) LIMIT ?"
	rows, err := repo.db.Query(selectQuery, message.Sender, message.Receiver, message.Receiver, message.Sender, 1)
	if err != nil {
		log.Println("failed to select messages %w", err)
		return fmt.Errorf("failed to select messages %w", err)
	}
	if rows.Next() {

		chatMembersQuery := `INSERT INTO chat_members ( user_id) VALUES (?)`

		result, err := repo.db.Exec(chatMembersQuery, message.Sender)
		if err != nil {
			log.Println("ffailed to insert chatMembersQuery %w", err)
			return fmt.Errorf("failed to insert chatMembersQuery %w", err)
		}
		// Access the result
		affectedRows, err := result.RowsAffected()
		if err != nil {
			log.Println("failed to get affected rows: %w", err)
			return fmt.Errorf("failed to get affected rows: %w", err)
		}
		fmt.Printf("Number of rows affected: %d\n", affectedRows)

		lastInsertID, err = result.LastInsertId()
		if err != nil {
			log.Println("failed to get last insert ID: %w", err)
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
	}
	query := `INSERT INTO messages (id,chat_id, sender, receiver, content, timestamp, is_forwarded, original_sender, original_message_id, is_edited, is_deleted)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ? ,?)`

	newUUID := uuid.New()

	uuidString := newUUID.String()

	message.ID = message.Sender + "|" + uuidString

	message.Sender = utils.RemoveAllButNumbersAndPlus(message.Sender)
	message.Receiver = utils.RemoveAllButNumbersAndPlus(message.Receiver)

	_, err = repo.db.Exec(query, message.ID, lastInsertID, message.Sender, message.Receiver, message.Content, message.Timestamp, message.IsForwarded, message.OriginalSender, message.OriginalMessageID, message.IsEdited, message.IsDeleted)

	if err != nil {
		log.Println("SaveMessage error ", err)
	}
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
	row := repo.db.QueryRow(`SELECT id, sender, receiver, content, timestamp, is_forwarded, original_sender, original_message_id, is_edited, is_deleted
	FROM messages WHERE id = ?`, id)

	err := row.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.Timestamp, &message.IsForwarded, &message.OriginalSender, &message.OriginalMessageID)
	if err == sql.ErrNoRows {
		return message, nil // No error, but no record found
	} else if err != nil {
		return message, err
	}
	return message, nil
}

// UpdateMessageContent updates the content of a message and sets it as edited
func (repo *MessageRepository) UpdateMessageContent(id, content string) error {
	_, err := repo.db.Exec(`UPDATE messages SET content = ?, is_edited = 1 WHERE id = ?`, content, id)
	return err
}

// DeleteMessage marks a message as deleted
func (repo *MessageRepository) DeleteMessage(id string) error {
	_, err := repo.db.Exec(`UPDATE messages SET is_deleted = 1 WHERE id = ?`, id)
	return err
}
