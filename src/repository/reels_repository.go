package repository

import (
	"database/sql"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type ReelRepository struct {
	DB *sql.DB
}

func (repo *ReelRepository) CreateReel(reel *models.Reel) error {
	query := "INSERT INTO reels (user_id, video_url, created_at) VALUES (?, ?, ?)"
	_, err := repo.DB.Exec(query, reel.UserID, reel.VideoURL, reel.CreatedAt)
	return err
}

func (repo *ReelRepository) FetchReels(lastID int) ([]models.Reel, error) {
	query := "SELECT id, user_id, video_url, created_at FROM reels WHERE id < ? ORDER BY created_at DESC LIMIT 10"
	rows, err := repo.DB.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reels := []models.Reel{}
	for rows.Next() {
		var reel models.Reel
		if err := rows.Scan(&reel.ID, &reel.UserID, &reel.VideoURL, &reel.CreatedAt); err != nil {
			return nil, err
		}
		reels = append(reels, reel)
	}
	return reels, nil
}

// Implement other methods like LikeReel, CommentReel, etc.

func (repo *ReelRepository) LikeReel(userID, reelID int) error {
	query := "INSERT INTO likes (user_id, reel_id) VALUES (?, ?)"
	_, err := repo.DB.Exec(query, userID, reelID)
	return err
}

func (repo *ReelRepository) UnlikeReel(userID, reelID int) error {
	query := "DELETE FROM likes WHERE user_id = ? AND reel_id = ?"
	_, err := repo.DB.Exec(query, userID, reelID)
	return err
}

func (repo *ReelRepository) CommentOnReel(comment *models.Comment) error {
	query := "INSERT INTO comments (user_id, reel_id, content, created_at) VALUES (?, ?, ?, ?)"
	_, err := repo.DB.Exec(query, comment.UserID, comment.ReelID, comment.Content, comment.CreatedAt)
	return err
}

func (repo *ReelRepository) GetCommentsForReel(reelID int) ([]models.Comment, error) {
	query := "SELECT id, user_id, reel_id, content, created_at FROM comments WHERE reel_id = ? ORDER BY created_at DESC"
	rows, err := repo.DB.Query(query, reelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.ReelID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
