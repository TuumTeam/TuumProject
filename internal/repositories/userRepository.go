package repositories

import (
	"database/sql"
	"tuum.com/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.CreatedAt)
	return err
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
