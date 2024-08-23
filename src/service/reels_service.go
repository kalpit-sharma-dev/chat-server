package service

import (
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
)

type ReelService struct {
	Repo *repository.ReelRepository
}

func (service *ReelService) CreateReel(reel *models.Reel) error {
	return service.Repo.CreateReel(reel)
}

func (service *ReelService) FetchReels(lastID int, limit int) ([]models.Reel, error) {
	return service.Repo.FetchReels(lastID, limit)
}

// Implement other methods like LikeReel, CommentReel, etc.

func (service *ReelService) LikeReel(userID, reelID int) error {
	return service.Repo.LikeReel(userID, reelID)
}

func (service *ReelService) UnlikeReel(userID, reelID int) error {
	return service.Repo.UnlikeReel(userID, reelID)
}

func (service *ReelService) CommentOnReel(comment *models.Comment) error {
	return service.Repo.CommentOnReel(comment)
}

func (service *ReelService) GetCommentsForReel(reelID int) ([]models.Comment, error) {
	return service.Repo.GetCommentsForReel(reelID)
}
