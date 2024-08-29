package repository

import (
	"database/sql"
	"fmt"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (username,phone_number,password_hash, verification_code,verified) VALUES (?,?,?,?,?)", user.UserName, user.Phone, user.Password, user.VerificationCode, true)
	return err
}

func (repo *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := repo.DB.QueryRow("SELECT id,username,phone_number,password_hash,IFNULL(full_name,''),IFNULL(profile_picture_url,''),IFNULL(status_message,''),IFNULL(verification_code,''),IFNULL(verified,0),IFNULL(created_at,''),IFNULL(updated_at,'') FROM users WHERE phone_number=?", phone).Scan(&user.ID, &user.UserName, &user.Phone, &user.Password, &user.FullName, &user.ProfilePictureUrl, &user.StatusMessage, &user.VerificationCode, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	_, err := repo.DB.Exec("UPDATE users SET verified=? WHERE phone=?", user.Verified, user.Phone)
	return err
}

func (repo *UserRepository) CheckUserInDB(phoneNumber string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = ?)"
	err := repo.DB.QueryRow(query, phoneNumber).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("could not execute query: %v", err)
	}

	return exists, nil
}
