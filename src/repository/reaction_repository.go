package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

// ReactionRepository provides CRUD operations for reactions
type ReactionRepository interface {
	Save(reaction models.Reaction) error
	GetReactionsByMessageID(messageID string) ([]models.Reaction, error)
}

type reactionRepository struct {
	db *sqlx.DB
}

func NewReactionRepository(db *sqlx.DB) ReactionRepository {
	return &reactionRepository{db: db}
}

// Save stores a reaction in the database
func (repo *reactionRepository) Save(reaction models.Reaction) error {
	_, err := repo.db.Exec(`INSERT INTO reactions (id, message_id, user, emoji, timestamp)
		VALUES (?, ?, ?, ?, ?)`,
		reaction.ID, reaction.MessageID, reaction.User, reaction.Emoji, reaction.Timestamp)
	return err
}

// GetReactionsByMessageID retrieves reactions for a specific message
func (repo *reactionRepository) GetReactionsByMessageID(messageID string) ([]models.Reaction, error) {
	rows, err := repo.db.Query(`SELECT id, message_id, user, emoji, timestamp FROM reactions WHERE message_id = ?`, messageID)
	if err != nil {
		return nil, err
	}
	// /defer rows.Close()

	var reactions []models.Reaction
	for rows.Next() {
		var reaction models.Reaction
		if err := rows.Scan(&reaction.ID, &reaction.MessageID, &reaction.User, &reaction.Emoji, &reaction.Timestamp); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
