package repository

import "github.com/kalpit-sharma-dev/chat-service/src/models"

type DatabaseRepository interface {
	CreateUser(user *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type IUserRepository interface {
	CreateUser(user *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
	UpdateUser(user *models.User) error
	CheckUserInDB(phoneNumber string) (bool, error)
}
