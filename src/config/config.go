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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT,
		receiver TEXT,
		content TEXT,
		timestamp TEXT,
		is_forwarded BOOLEAN,
		original_sender TEXT,
		original_message_id TEXT

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

	return db
}
