package repositories

import (
	"database/sql"

	"tuum.com/internal/models"
	"tuum.com/internal/services"
)

// UserRepository définit les méthodes pour interagir avec la table des utilisateurs.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository crée un nouveau UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create insère un nouvel utilisateur dans la base de données.
func (repo *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		return err
	}

	stmt, err := repo.DB.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, hashedPassword)
	return err
}

// FindByEmail recherche un utilisateur par email.
func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	row := repo.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
