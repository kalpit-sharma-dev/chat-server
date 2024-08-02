package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/service"
	"github.com/kalpit-sharma-dev/chat-service/src/utils"
)

type UserController struct {
	UserService service.IUserService
}

func NewUserController(userService service.IUserService) UserController {
	return UserController{UserService: userService}
}

func (controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.User
	phone := r.FormValue("phone")
	userName := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	req.Email = email
	req.UserName = userName
	req.Password = password
	req.Phone = phone
	err := controller.UserService.RegisterUser(req, phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful, verification code sent"})
}

func (controller *UserController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phone")
	code := r.FormValue("code")

	err := controller.UserService.VerifyUser(phone, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification successful"})
}

func (controller *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phone")

	err := controller.UserService.LoginUser(phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(phone)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": token})
}
