package repositories

import (
	"database/sql"

	"tuum.com/internal/models"
)

// UserRepository définit les méthodes pour interagir avec la table des post.
type PostRepository struct {
	DB *sql.DB
}

// NewUPostRepository crée un nouveau PostRepository.
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (repo *PostRepository) CreatePost(post *models.Post) error {

	stmt, err := repo.DB.Prepare("INSERT INTO post (Title, Content, user, Comments) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.ID, post.UserID, post.RoomID, post.Title, post.Content, post.CreatedAt)
	return err
}
