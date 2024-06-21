package models

type Salon struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required,min=3,max=32"`
	Description string `json:"description" validate:"required,min=3,max=500"`
	Post        []Post `json:"post" validate:"required,email"`
}
