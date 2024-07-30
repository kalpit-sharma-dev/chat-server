package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
)

type Client struct {
	conn  *websocket.Conn
	phone string
}

type ChatService struct {
	clients        map[string]*Client
	groupClients   map[int64][]*Client
	broadcast      chan models.Message
	groupBroadcast chan models.GroupMessage
	messageRepo    *repository.MessageRepository
	groupRepo      *repository.GroupRepository
	mu             sync.RWMutex
}

func NewChatService(messageRepo *repository.MessageRepository, groupRepo *repository.GroupRepository) *ChatService {
	return &ChatService{
		clients:        make(map[string]*Client),
		groupClients:   make(map[int64][]*Client),
		broadcast:      make(chan models.Message),
		groupBroadcast: make(chan models.GroupMessage),
		messageRepo:    messageRepo,
		groupRepo:      groupRepo,
	}
}

func (service *ChatService) AddClient(phone string, conn *websocket.Conn) {
	service.mu.Lock()
	service.clients[phone] = &Client{conn: conn, phone: phone}
	service.mu.Unlock()
}

func (service *ChatService) RemoveClient(phone string) {
	service.mu.Lock()
	delete(service.clients, phone)
	service.mu.Unlock()
}

func (service *ChatService) BroadcastMessage(message models.Message) {
	service.broadcast <- message
}

func (service *ChatService) BroadcastGroupMessage(message models.GroupMessage) {
	service.groupBroadcast <- message
}

func (service *ChatService) HandleMessages() {
	for {
		select {
		case message := <-service.broadcast:
			service.messageRepo.SaveMessage(&message)
			service.mu.Lock()
			if client, ok := service.clients[message.Receiver]; ok {
				err := client.conn.WriteJSON(message)
				if err != nil {
					client.conn.Close()
					delete(service.clients, message.Receiver)
				}
			}
			service.mu.Unlock()
		case groupMessage := <-service.groupBroadcast:
			service.groupRepo.SaveGroupMessage(&groupMessage)
			members, err := service.groupRepo.GetGroupMembers(groupMessage.GroupID)
			if err != nil {
				continue
			}
			service.mu.Lock()
			for _, member := range members {
				if client, ok := service.clients[member]; ok {
					err := client.conn.WriteJSON(groupMessage)
					if err != nil {
						client.conn.Close()
						delete(service.clients, member)
					}
				}
			}
			service.mu.Unlock()
		}
	}
}

func (service *ChatService) GetChatHistory(sender, receiver string) ([]models.Message, error) {
	return service.messageRepo.GetMessages(sender, receiver)
}

func (service *ChatService) GetGroupChatHistory(groupID int64) ([]models.GroupMessage, error) {
	return service.groupRepo.GetGroupMessages(groupID)
}

func (service *ChatService) CreateGroup(group *models.Group) error {
	return service.groupRepo.CreateGroup(group)
}

// GetClient retrieves the WebSocket connection for a user
func (service *ChatService) GetClient(phoneNumber string) *websocket.Conn {
	service.mu.RLock()
	if client, ok := service.clients[phoneNumber]; ok {
		return client.conn
	}
	return nil

}

// ForwardMessage forwards a message to a new receiver or group
func (service *ChatService) ForwardMessage(originalMessageID, newReceiver string) error {
	// Retrieve the original message
	originalMessage, err := service.messageRepo.GetMessageByID(originalMessageID)
	if err != nil {
		return err
	}

	// Create a new message for forwarding
	forwardedMessage := models.Message{
		ID:                generateNewID(), // Generate a new ID for the forwarded message
		Sender:            originalMessage.Sender,
		Receiver:          newReceiver,
		Content:           originalMessage.Content,
		Timestamp:         time.Now().Format(time.RFC3339),
		IsForwarded:       true,
		OriginalSender:    originalMessage.Sender,
		OriginalMessageID: originalMessageID,
	}

	// Save the forwarded message
	if err := service.messageRepo.SaveMessage(&forwardedMessage); err != nil {
		return err
	}

	// Broadcast the forwarded message
	if client := service.GetClient(newReceiver); client != nil {
		client.WriteJSON(forwardedMessage)
		//client.conn.WriteJSON(forwardedMessage)
	}

	return nil
}

// generateNewID generates a new unique identifier for messages
func generateNewID() string {
	return uuid.NewString()
}
