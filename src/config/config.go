package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "./project.db")
	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		sender TEXT,
		receiver TEXT,
		receiver_id INTEGER,
		chat_id INTEGER,
		sender_id INTEGER,
		content TEXT,
		timestamp TEXT,
		is_forwarded BOOLEAN,
		original_sender TEXT,
		original_message_id TEXT,
		is_edited BOOLEAN DEFAULT 0,
		is_deleted BOOLEAN DEFAULT 0,
		FOREIGN KEY (chat_id) REFERENCES chats (id),
    	FOREIGN KEY (sender_id) REFERENCES users (id)
	);
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	CREATE TABLE IF NOT EXISTS group_members (
		group_id INTEGER,
		member TEXT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);
	CREATE TABLE IF NOT EXISTS group_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER,
		sender TEXT,
		content TEXT,
		timestamp TEXT,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);
	CREATE TABLE IF NOT EXISTS media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		type TEXT,
		message_id INTEGER,
		FOREIGN KEY (message_id) REFERENCES messages(id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
	// Create reactions table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reactions (
		id TEXT PRIMARY KEY,
		message_id TEXT,
		user TEXT,
		emoji TEXT,
		timestamp TEXT,
		FOREIGN KEY(message_id) REFERENCES messages(id)
	);
	
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phone_number TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    is_group BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE IF NOT EXISTS chat_members (
    chat_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT,
    FOREIGN KEY (chat_id) REFERENCES chats (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
	UNIQUE (chat_id, user_id)

);
	`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
