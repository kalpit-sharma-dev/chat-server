package main

import (
	"log"
	"runtime/debug"

	_ "modernc.org/sqlite"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kalpit-sharma-dev/chat-service/src/server"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Print(r, debug.Stack())
		}
	}()
	log.Print("INFO : starting the application ....")
	log.Print("INFO : chat-service started ....")

	server.LoadRoute()
}

// func initDB() *sql.DB {
// 	db, err := sql.Open("sqlite", "file:chat.db?cache=shared&mode=rwc")
// 	if err != nil {
// 		log.Fatalf("Failed to open database: %v", err)
// 	}

// 	createTable := `CREATE TABLE IF NOT EXISTS messages (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		room TEXT,
// 		sender TEXT,
// 		content TEXT,
// 		private BOOLEAN,
// 		recipient TEXT,
// 		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 	);`

// 	_, err = db.Exec(createTable)
// 	if err != nil {
// 		log.Fatalf("Failed to create table: %v", err)
// 	}

// 	return db
// }

// func saveMessage(db *sql.DB, msg Message) {
// 	_, err := db.Exec(`INSERT INTO messages (room, sender, content, private, recipient, time) VALUES (?, ?, ?, ?, ?, ?)`,
// 		msg.Room, msg.Sender, msg.Content, msg.Private, msg.To, msg.Time)
// 	if err != nil {
// 		log.Printf("Failed to save message: %v", err)
// 	}
// }

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.clients[client] = true
// 		case client := <-h.unregister:
// 			if _, ok := h.clients[client]; ok {
// 				delete(h.clients, client)
// 				close(client.send)
// 				h.mutex.Lock()
// 				delete(h.privateMsg, client.username)
// 				h.mutex.Unlock()
// 				for _, room := range h.rooms {
// 					delete(room.clients, client)
// 				}
// 			}
// 		case message := <-h.broadcast:
// 			if message.Private {
// 				if recipient, ok := h.privateMsg[message.To]; ok {
// 					recipient.send <- []byte(message.Content)
// 					saveMessage(h.db, message)
// 				}
// 			} else if room, ok := h.rooms[message.Room]; ok {
// 				for client := range room.clients {
// 					select {
// 					case client.send <- []byte(fmt.Sprintf("%s: %s", message.Sender, message.Content)):
// 					default:
// 						close(client.send)
// 						delete(room.clients, client)
// 					}
// 				}
// 				saveMessage(h.db, message)
// 			}
// 		}
// 	}
// }

// func main() {
// 	db := initDB()
// 	defer db.Close()

// 	hub := NewHub(db)
// 	go hub.Run()
// 	r := mux.NewRouter()
// 	http.HandleFunc("/login", login)
// 	r.HandleFunc("/register", registerUser).Methods("POST")
// 	r.HandleFunc("/verify", verifyUser).Methods("POST")
// 	http.HandleFunc("/upload", uploadFile)
// 	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files

// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		hub.ServeWs(w, r)
// 	})

// 	fmt.Println("Chat server started on :9999")
// 	log.Fatal(http.ListenAndServe(":9999", nil))
// }

// // generateVerificationCode generates a random 6-digit verification code
// func generateVerificationCode() string {
// 	rand.Seed(uint64(time.Now().UnixNano()))
// 	return fmt.Sprintf("%06d", rand.Intn(1000000))
// }

// // sendVerificationCode sends an SMS with the verification code using Twilio
// func sendVerificationCode(phone string, code string) error {
// 	client := twilio.NewRestClient()

// 	params := &openapi.CreateMessageParams{}
// 	params.SetTo(phone)
// 	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
// 	params.SetBody(fmt.Sprintf("Your verification code is: %s", code))

// 	_, err := client.Api.CreateMessage(params)
// 	return err
// }

// var clientOptions = options.Client().ApplyURI(os.Getenv("MONGO_URI"))
// var client, err = mongo.Connect(context.TODO(), clientOptions)

// //	if err != nil {
// //		log.Fatal(err)
// //	}
// //
// // err = client.Ping(context.TODO(), nil)
// //
// //	if err != nil {
// //		log.Fatal(err)
// //	}
// func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Failed to upgrade to WebSocket:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		var message models.Message
// 		err := conn.ReadJSON(&message)
// 		if err != nil {
// 			log.Println("Error reading message:", err)
// 			break
// 		}

// 		message.Timestamp = time.Now()
// 		collection := client.Database("whatsapp").Collection("messages")
// 		_, err = collection.InsertOne(context.TODO(), message)
// 		if err != nil {
// 			log.Println("Error inserting message:", err)
// 			continue
// 		}

// 		// Echo the message back to the client
// 		err = conn.WriteJSON(message)
// 		if err != nil {
// 			log.Println("Error sending message:", err)
// 			break
// 		}
// 	}
// }
