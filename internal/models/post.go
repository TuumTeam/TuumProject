package models

type Post struct {
	Title    string    `json:"title" validate:"required,min=3,max=32"`
	Content  string    `json:"content" validate:"required"`
	Hashtag  string    `json:"hashtag" omitempty:"null"`
	User     User      `json:"user" validate:"required"`
	Comments []Comment `json:"comments" validate:"required"`
}
