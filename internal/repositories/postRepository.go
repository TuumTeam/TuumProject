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
func NewPostRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreatePost(post *models.Post) error {

	stmt, err := repo.DB.Prepare("INSERT INTO users (Title, Content, user, Comments) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Title, post.Content, post.User, post.Comments)
	return err
}

// FindByTitle recherche un post par Title.
/*func (repo *PostRepository) FindByTitle(title string) (*models.Post, error) {
	post := &models.Post{}
	row := repo.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", title)
	err := row.Scan(&post.Title, &post.Content, &post.User, &post.Comments)
	if err != nil {
		return nil, err
	}
	return post, nil
}*/
