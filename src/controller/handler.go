package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/service"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/exp/rand"
)

var jwtKey = []byte("your_secret_key")

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

type Handler struct {
	service service.ChatService
	Hub     models.Hub
	C       *models.Client
}

func NewHandler(service service.ChatService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

// registerUser handles the user registration process
// func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
// 	phone := r.FormValue("phone")

// 	// Check if the user already exists
// 	var exists bool
// 	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone=?)", phone).Scan(&exists)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if exists {
// 		http.Error(w, "User already exists", http.StatusConflict)
// 		return
// 	}

// 	// Generate and send verification code
// 	verificationCode := generateVerificationCode()
// 	err = sendVerificationCode(phone, verificationCode)
// 	if err != nil {
// 		http.Error(w, "Failed to send verification code", http.StatusInternalServerError)
// 		return
// 	}

// 	// Store user in the database with verification code (for simplicity, not storing in this example)
// 	_, err = db.Exec("INSERT INTO users (phone, verification_code) VALUES (?, ?)", phone, verificationCode)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful, verification code sent"})
// }

// generateVerificationCode generates a random 6-digit verification code
func generateVerificationCode() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// sendVerificationCode sends an SMS with the verification code using Twilio
func sendVerificationCode(phone string, code string) error {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(fmt.Sprintf("Your verification code is: %s", code))

	_, err := client.Api.CreateMessage(params)
	return err
}

// // verifyUser handles the verification of the user using the verification code
// func (h *Handler) VerifyUser(w http.ResponseWriter, r *http.Request) {
// 	phone := r.FormValue("phone")
// 	code := r.FormValue("code")

// 	// Verify the code
// 	var storedCode string
// 	err := db.QueryRow("SELECT verification_code FROM users WHERE phone=?", phone).Scan(&storedCode)
// 	if err != nil {
// 		http.Error(w, "Invalid phone number or verification code", http.StatusUnauthorized)
// 		return
// 	}
// 	if storedCode != code {
// 		http.Error(w, "Invalid verification code", http.StatusUnauthorized)
// 		return
// 	}

// 	// Update user as verified (for simplicity, not updating in this example)
// 	_, err = db.Exec("UPDATE users SET verified=1 WHERE phone=?", phone)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Verification successful"})
// }

// func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
// 	file, _, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	out, err := os.Create("/path/to/save/file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Write([]byte("File uploaded successfully"))
// }

// func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
// 	tokenStr := r.URL.Query().Get("token")
// 	claims := &models.Claims{}
// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	client := &models.Client{Username: claims.Username, Conn: conn, Send: make(chan []byte, 256)}
// 	h.Hub.Register <- client
// 	h.C = client

// 	h.Hub.Mutex.Lock()
// 	h.Hub.PrivateMsg[client.Username] = client
// 	h.Hub.Mutex.Unlock()

// 	go h.ReadPump()
// 	go h.WritePump()
// }

// func (h *Handler) ReadPump() {
// 	defer func() {
// 		h.Hub.Unregister <- h.C
// 		h.C.Conn.Close()
// 	}()

// 	for {
// 		_, message, err := h.C.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("Unexpected close error: %v", err)
// 			}
// 			break
// 		}

// 		msg := models.Message{}
// 		err = json.Unmarshal(message, &msg)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		msg.Time = time.Now()
// 		msg.Sender = h.C.Username

// 		if msg.Private {
// 			h.Hub.Broadcast <- msg
// 		} else {
// 			h.Hub.Mutex.Lock()
// 			if room, ok := h.Hub.Rooms[msg.Room]; ok {
// 				room.Clients[h.C] = true
// 			} else {
// 				h.Hub.Rooms[msg.Room] = &models.Room{Name: msg.Room, Clients: map[*models.Client]bool{h.C: true}}
// 			}
// 			h.Hub.Mutex.Unlock()
// 			h.Hub.Broadcast <- msg
// 		}
// 	}
// }

// func (h *Handler) WritePump() {
// 	defer func() {
// 		h.C.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-h.C.Send:
// 			if !ok {
// 				h.C.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			if err := h.C.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
