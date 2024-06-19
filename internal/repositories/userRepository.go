package repositories

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"tuum.com/internal/models"
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
func (repo *UserRepository) Create(user *models.User) error {
	hashedPassword, err := hashPassword(user.Password)
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

// hashPassword génère un hash pour un mot de passe donné.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash vérifie si un mot de passe correspond à un hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
