package service

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
)

type UserService struct {
	UserRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &UserService{UserRepo: userRepo}
}

func (service *UserService) GenerateVerificationCode() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

func (service *UserService) SendVerificationCode(phone string, code string) error {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(fmt.Sprintf("Your verification code is: %s", code))

	_, err := client.Api.CreateMessage(params)
	return err
}

func (service *UserService) RegisterUser(user models.User, phone string) error {
	existingUser, _ := service.UserRepo.GetUserByPhone(phone)
	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	verificationCode := service.GenerateVerificationCode()
	err := service.SendVerificationCode(phone, verificationCode)
	if err == nil {
		return err
	}
	user.VerificationCode = verificationCode
	// user := &models.User{
	// 	Phone:            phone,
	// 	VerificationCode: verificationCode,
	// }

	return service.UserRepo.CreateUser(&user)
}

func (service *UserService) VerifyUser(phone string, code string) error {
	user, err := service.UserRepo.GetUserByPhone(phone)
	if err != nil {
		return fmt.Errorf("invalid phone number or verification code")
	}
	if user.VerificationCode != code {
		return fmt.Errorf("invalid verification code")
	}

	user.Verified = true
	return service.UserRepo.UpdateUser(user)
}

func (service *UserService) LoginUser(phone string) error {
	user, err := service.UserRepo.GetUserByPhone(phone)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if !user.Verified {
		return fmt.Errorf("user not verified")
	}
	return nil
}
