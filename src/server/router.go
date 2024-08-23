package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kalpit-sharma-dev/chat-service/src/config"
	"github.com/kalpit-sharma-dev/chat-service/src/controller"
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
	"github.com/kalpit-sharma-dev/chat-service/src/service"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func LoadRoute() {
	log.Println("INFO : Loading Router")
	router := mux.NewRouter().PathPrefix("/chat-service/api").Subrouter()
	router.Use(headerMiddleware)
	registerAppRoutes(router)
	log.Println("INFO : Router Loaded Successfully")
	log.Println("INFO : Application is started Successfully")
	// wg := &sync.WaitGroup{}
	// go LoadSlots(wg)
	http.ListenAndServe(":9999", router)
	//*****secure TLS********log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
}

var SlotCounter uint64

func registerAppRoutes(r *mux.Router) {
	log.Println("INFO : Registering Router ")

	log.Println("INFO : Registering Router ")

	var err error
	// Connect to MySQL database
	// dbmySqlCon, err := sql.Open("mysql", "kalpit:password@tcp(192.168.100.4:3306)/demo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dbmySqlCon, err := sql.Open("mysql", "kalpit:password@tcp(localhost:3306)/chatserver")
	if err != nil {
		log.Fatal(err)
	}
	// Test database connection
	if err := dbmySqlCon.Ping(); err != nil {
		log.Fatal(err)
	}
	db := config.InitDB()
	//var dbConn db.DatabaseImpl

	userRepo := repository.NewUserRepository(dbmySqlCon)
	// /eventRepo := db.NewDatabaseRepository(dbmySqlCon)

	//userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	//defer db.Close()

	messageRepo := repository.NewMessageRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	chatRepo := repository.NewChatRepository(db)
	reactionRepo := repository.NewReactionRepository(db)

	chatService := service.NewChatService(messageRepo, groupRepo, reactionRepo, chatRepo)
	awsSession := session.Must(session.NewSession())
	s3Client := s3.New(awsSession)
	mediaService := service.NewMediaService(mediaRepo, s3Client, "your-s3-bucket-name")
	chatController := controller.NewChatController(chatService, mediaService)

	//reels
	reelRepository := &repository.ReelRepository{DB: db.DB}
	reelService := &service.ReelService{Repo: reelRepository}
	reelController := &controller.ReelController{Service: reelService}

	go chatService.HandleMessages()

	r.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	r.HandleFunc("/verify", userController.VerifyUser).Methods("POST")
	r.HandleFunc("/login", userController.LoginUser).Methods("POST")

	r.HandleFunc("/ws", chatController.HandleWebSocket)

	r.HandleFunc("/groups", chatController.CreateGroup).Methods("POST")
	r.HandleFunc("/group_messages", chatController.GetGroupMessages).Methods("GET")
	r.HandleFunc("/reactions", chatController.AddReactionHandler(chatService)).Methods("POST")

	r.HandleFunc("/upload_media", chatController.UploadMedia).Methods("POST")

	r.HandleFunc("/messages/edit", chatController.EditMessageHandler(chatService)).Methods("POST")
	r.HandleFunc("/messages/delete", chatController.DeleteMessageHandler(chatService)).Methods("POST")
	r.HandleFunc("/chats/{phone}", chatController.GetChats).Methods("GET")

	//reels
	r.HandleFunc("/reels/upload", reelController.UploadReel).Methods("POST")
	r.HandleFunc("/reels", reelController.FetchReels).Methods("GET")
	r.HandleFunc("/reels/{id}/like", reelController.LikeReel).Methods("POST")
	r.HandleFunc("/reels/{id}/unlike", reelController.UnlikeReel).Methods("POST")
	r.HandleFunc("/reels/{id}/comment", reelController.CommentOnReel).Methods("POST")
	r.HandleFunc("/reels/{id}/comments", reelController.GetCommentsForReel).Methods("GET")

	//r.HandleFunc("/login", chatHandlers.Login).Methods(http.MethodPost)

	//r.HandleFunc("/register", chatHandlers.RegisterUser).Methods(http.MethodGet) //quer param color,number

	//r.HandleFunc("/verify", chatHandlers.VerifyUser).Methods(http.MethodGet) ////quer param color return slot numbers

	//r.HandleFunc("/upload", chatHandlers.UploadFile).Methods(http.MethodPost)

	// r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	chatHandlers.ServeWs(w, r)
	// }).Methods(http.MethodGet)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("."))
	}).Methods(http.MethodGet)

	//go hub.Run()
	log.Println("INFO : Router Registered Successfully")
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
