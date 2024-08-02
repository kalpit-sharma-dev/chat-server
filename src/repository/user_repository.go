package repository

import (
	"database/sql"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (user_name,phone_number, verification_code) VALUES (?, ?)", user.UserName, user.Phone, user.VerificationCode)
	return err
}

func (repo *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := repo.DB.QueryRow("SELECT id, phone, verification_code, verified, created_at FROM users WHERE phone=?", phone).Scan(&user.ID, &user.Phone, &user.VerificationCode, &user.Verified, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	_, err := repo.DB.Exec("UPDATE users SET verified=? WHERE phone=?", user.Verified, user.Phone)
	return err
}
