package repository

import (
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
func (repo *ChatRepository) GetChatsForUser(userID int) ([]models.Chat, error) {
	rows, err := repo.db.Query(`
        SELECT c.id, c.name, c.is_group
        FROM chats c
        JOIN chat_members cm ON c.id = cm.chat_id
        WHERE cm.user_id = ?
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsGroup); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
