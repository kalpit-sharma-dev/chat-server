package models

import "time"

type Reel struct {
	ID        int
	UserID    int
	VideoURL  string
	CreatedAt time.Time
}

type Like struct {
	ID     int
	UserID int
	ReelID int
}

type Comment struct {
	ID        int
	UserID    int
	ReelID    int
	Content   string
	CreatedAt time.Time
}
