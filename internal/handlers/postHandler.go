package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"tuum.com/internal/models"
	"tuum.com/internal/services"

	"github.com/gorilla/mux"
)

type PostHandler struct {
	PostService services.PostService
}

func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{PostService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.CreatedAt = time.Now()

	if err := h.PostService.CreatePost(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	posts, err := h.PostService.GetPostsByUserID(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		http.Error(w, "Unable to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
