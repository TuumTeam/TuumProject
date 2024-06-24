package models

import "time"

type Post struct {
	ID        int
	UserID    int
	RoomID    int
	Title     string
	Content   string
	CreatedAt time.Time
}
