package models

type Comment struct {
	Content string `json:"content" validate:"required"`
	User    User   `json:"user" validate:"required"`
}
