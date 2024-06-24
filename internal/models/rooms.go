package models

import "time"

type Room struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}
