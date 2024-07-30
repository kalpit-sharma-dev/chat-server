package service

// type ChatService interface {
// 	RegisterUser()
// 	Login()
// 	VerifyUser()
// 	UploadFile()
// 	ServeWs()
// }

type IUserService interface {
	RegisterUser(phone string) error
	SendVerificationCode(phone string, code string) error
	GenerateVerificationCode() string
	VerifyUser(phone string, code string) error
	LoginUser(phone string) error
}
