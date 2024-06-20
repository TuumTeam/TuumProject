package handlers

import (
	"encoding/json"
	"net/http"
	"tuum.com/internal/models"
	"tuum.com/internal/repositories"
)

// UserHandler gère les requêtes liées aux utilisateurs.
type UserHandler struct {
	Repo *repositories.UserRepository
}

// NewUserHandler crée un nouveau UserHandler.
func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// Signup gère l'inscription des nouveaux utilisateurs.
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Repo.Create(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login gère la connexion des utilisateurs.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.FindByEmail(credentials.Email)
	if err != nil || !repositories.CheckPasswordHash(credentials.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Générer et renvoyer un token JWT (non implémenté ici)
	w.WriteHeader(http.StatusOK)
}
