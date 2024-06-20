package services

import (
	"tuum.com/internal/models"
	"tuum.com/internal/repositories"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	GetPostsByUserID(userID int) ([]*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
}

type postService struct {
	postRepository repositories.PostRepository
}

func NewPostService(postRepo repositories.PostRepository) PostService {
	return &postService{
		postRepository: postRepo,
	}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.postRepository.CreatePost(post)
}

func (s *postService) GetPostByID(id int) (*models.Post, error) {
	return s.postRepository.GetPostByID(id)
}

func (s *postService) GetPostsByUserID(userID int) ([]*models.Post, error) {
	return s.postRepository.GetPostsByUserID(userID)
}

func (s *postService) GetAllPosts() ([]*models.Post, error) {
	return s.postRepository.GetAllPosts()
}
