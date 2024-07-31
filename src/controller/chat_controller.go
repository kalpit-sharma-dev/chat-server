package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/service"
	"github.com/kalpit-sharma-dev/chat-service/src/utils"
)

type ChatController struct {
	ChatService  *service.ChatService
	MediaService *service.MediaService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatController(chatService *service.ChatService, mediaService *service.MediaService) *ChatController {
	return &ChatController{
		ChatService:  chatService,
		MediaService: mediaService,
	}
}

func (controller *ChatController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	controller.ChatService.AddClient(claims.Phone, conn)
	defer controller.ChatService.RemoveClient(claims.Phone)

	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		if groupID, ok := message["group_id"].(float64); ok {
			groupMessage := models.GroupMessage{
				GroupID:   int64(groupID),
				Sender:    claims.Phone,
				Content:   message["content"].(string),
				Timestamp: time.Now().Format(time.RFC3339),
			}
			controller.ChatService.BroadcastGroupMessage(groupMessage)
		} else if messageType, ok := message["type"].(string); ok {
			if messageType == "webrtc" {
				// Handle WebRTC signaling messages
				controller.handleSignalingMessage(message, conn)
			} else {
				chatMessage := models.Message{
					Sender:    claims.Phone,
					Receiver:  message["receiver"].(string),
					Content:   message["content"].(string),
					Timestamp: time.Now().Format(time.RFC3339),
				}
				controller.ChatService.BroadcastMessage(chatMessage)
			}
		}
	}
}
func (controller *ChatController) handleSignalingMessage(message map[string]interface{}, conn *websocket.Conn) {
	// Handle WebRTC signaling messages (offer, answer, ICE candidates)
	// Broadcast the signaling message to the intended recipient(s)
	recipient := message["recipient"].(string)
	conn = controller.ChatService.GetClient(recipient)
	if conn != nil {
		conn.WriteJSON(message)
	}
}
func (controller *ChatController) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = controller.ChatService.CreateGroup(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (controller *ChatController) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupIDStr := r.URL.Query().Get("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	messages, err := controller.ChatService.GetGroupChatHistory(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

func (controller *ChatController) UploadMedia(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := controller.MediaService.UploadFile(file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	media := models.Media{
		URL:  url,
		Type: r.FormValue("type"),
	}
	controller.MediaService.SaveMedia(&media)

	utils.RespondJSON(w, http.StatusCreated, media)
}

// AddReactionHandler handles the addition of reactions to messages
func (controller *ChatController) AddReactionHandler(chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AddReactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := chatService.AddReaction(req.MessageID, req.User, req.Emoji)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// EditMessageHandler handles editing a message
func (controller *ChatController) EditMessageHandler(chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.EditMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := chatService.EditMessage(req.MessageID, req.NewContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// DeleteMessageHandler handles deleting a message
func (controller *ChatController) DeleteMessageHandler(chatService *service.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.DeleteMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := chatService.DeleteMessage(req.MessageID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
