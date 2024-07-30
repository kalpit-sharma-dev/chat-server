package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type MediaRepository struct {
	db *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (repo *MediaRepository) SaveMedia(media *models.Media) error {
	query := `INSERT INTO media (url, type, message_id) VALUES (?, ?, ?)`
	_, err := repo.db.Exec(query, media.URL, media.Type, media.MessageID)
	return err
}

func (repo *MediaRepository) GetMediaByMessageID(messageID int64) ([]models.Media, error) {
	var media []models.Media
	query := `SELECT * FROM media WHERE message_id = ?`
	err := repo.db.Select(&media, query, messageID)
	return media, err
}
