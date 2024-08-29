package service

import "github.com/kalpit-sharma-dev/chat-service/src/models"

// type ChatService interface {
// 	RegisterUser()
// 	Login()
// 	VerifyUser()
// 	UploadFile()
// 	ServeWs()
// }

type IUserService interface {
	RegisterUser(user models.User, phone string) error
	SendVerificationCode(phone string, code string) error
	GenerateVerificationCode() string
	VerifyUser(phone string, code string) error
	LoginUser(phone string) error
	CheckUserService(phoneNumber string) (bool, error)
}
