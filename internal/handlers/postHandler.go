package handlers

import (
	"tuum.com/internal/repositories"
)

type PostHandler struct {
	Repo *repositories.PostRepository
}
