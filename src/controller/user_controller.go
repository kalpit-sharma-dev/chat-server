package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
	if len(phone) != 10 {
		http.Error(w, "phone number should be 10 digits", http.StatusBadRequest)
		return
	} else {
		phone = "+91" + phone
		req.Phone = phone

	}

	err := controller.UserService.RegisterUser(req, phone)
	if err != nil {
		http.Error(w, "user already registered", http.StatusConflict)
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

	password := r.FormValue("password")
	log.Println(phone)
	phone = utils.RemoveAllButNumbersAndPlus(phone)
	if len(phone) != 10 {
		http.Error(w, "phone number should be 10 digits", http.StatusBadRequest)
		return
	} else {
		phone = strings.ReplaceAll(phone, " ", "")
		phone = "+91" + phone
		log.Println(phone)
	}

	user := models.User{
		Phone:    phone,
		Password: password,
	}
	err := controller.UserService.LoginUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	log.Println("JWTToken for ", phone)
	token, err := utils.GenerateJWT(phone)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": token})
}

func (controller *UserController) CheckUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PhoneNumber string `json:"phoneNumber"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	phone := utils.RemoveAllButNumbersAndPlus(input.PhoneNumber)
	phone = strings.ReplaceAll(phone, " ", "")
	if contains := strings.Contains(phone, "+91"); !contains {
		phone = "+91" + phone
	}
	log.Println(phone)
	isRegistered, err := controller.UserService.CheckUserService(phone)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := models.CheckUserResposne{
		IsRegistered: isRegistered,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
