package models

import (
	"database/sql"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Message1 struct {
	ID        int       `json:"id"`
	Room      string    `json:"room,omitempty"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Private   bool      `json:"private"`
	To        string    `json:"to,omitempty"`
	Time      time.Time `json:"time"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	ID                string `json:"id"`
	Sender            string `json:"sender"`
	Receiver          string `json:"receiver"`
	Content           string `json:"content"` // This will include formatted text, emojis, and links
	Timestamp         string `json:"timestamp"`
	IsForwarded       bool   `json:"is_forwarded"`                  // Indicates if the message is forwarded
	OriginalSender    string `json:"original_sender,omitempty"`     // Original sender if forwarded
	OriginalMessageID string `json:"original_message_id,omitempty"` // ID of the original message if forwarded
}

type Client struct {
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

type Room struct {
	Name    string
	Clients map[*Client]bool
}

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string]*Room
	PrivateMsg map[string]*Client
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.Mutex
	Db         *sql.DB
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]*Room),
		PrivateMsg: make(map[string]*Client),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Db:         db,
	}
}

type User struct {
	ID               int
	Phone            string
	VerificationCode string
	Verified         bool
	CreatedAt        string
}

type Group struct {
	ID      int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type GroupMessage struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID   int64  `json:"group_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type Media struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	URL       string `json:"url"`
	Type      string `json:"type"` // e.g., image, video, document
	MessageID int64  `json:"message_id"`
}
