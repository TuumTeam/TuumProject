package repositories

import (
	"database/sql"
	"tuum.com/internal/models"
)

type PostRepository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id int) (*models.Post, error)
	GetPostsByUserID(userID int) ([]*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	query := `INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.UserID, post.Title, post.Content, post.CreatedAt)
	return err
}

func (r *postRepository) GetPostByID(id int) (*models.Post, error) {
	query := `SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var post models.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetPostsByUserID(userID int) ([]*models.Post, error) {
	query := `SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *postRepository) GetAllPosts() ([]*models.Post, error) {
	query := `SELECT id, user_id, title, content, created_at FROM posts`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}
