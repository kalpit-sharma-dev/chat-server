package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

// ChatRepository handles database operations related to chats.
type ChatRepository struct {
	db *sqlx.DB
}

// NewChatRepository creates a new ChatRepository.
func NewChatRepository(db *sqlx.DB) ChatRepository {
	return ChatRepository{db: db}
}

// GetChatsForUser retrieves the list of chats for a given user.
func (repo *ChatRepository) GetChatsForUser(userID, otherUser string) ([]models.Message, error) {

	log.Println(userID, otherUser)
	rows, err := repo.db.Query(`
        SELECT id, sender, receiver, content, timestamp, is_forwarded, original_sender, original_message_id, is_edited, is_deleted FROM messages m
        JOIN chat_members cm ON m.chat_id = cm.chat_id
        WHERE cm.user_id = ? and m.sender = ? and m.receiver = ?
    `, userID, userID, otherUser)
	if err != nil {
		log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		log.Println("error getting messages for a user", err)
		return nil, err
	}
	//defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.Timestamp, &message.IsForwarded, &message.OriginalSender, &message.OriginalMessageID, &message.IsEdited, &message.IsDeleted); err != nil {
			log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@", err)
			return nil, err
		}
		messages = append(messages, message)

	}

	return messages, nil
}
