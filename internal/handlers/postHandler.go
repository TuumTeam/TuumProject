package handlers

import (
	"tuum.com/internal/repositories"
)

type PostHandler struct {
	Repo *repositories.PostRepository
}

func NewPostHandler(repo *repositories.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}
