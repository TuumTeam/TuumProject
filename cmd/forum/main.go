package main

import (
	"log"
	"net/http"
	"time"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"tuum.com/internal/repositories"
	"tuum.com/internal/services"
	"tuum.com/internal/tuum_database"
	"tuum.com/pkg/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig("internal/config/config.yaml")

	// Initialize the database connection
	db, err := tuum_database.InitDB(&config.AppConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	tuum_database.RunMigrations(db, "./migrations")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// Set up router
	r := mux.NewRouter()

	r.Handle("/static/{path}", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Public routes
	r.HandleFunc("/", handlers.Main)
	r.HandleFunc("/login", handlers.Login)

	// Create a subrouter for authenticated routes
	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)

	// User routes
	s.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	s.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	s.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")

	// Post routes
	s.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	s.HandleFunc("/posts/{id}", postHandler.GetPostByID).Methods("GET")
	s.HandleFunc("/posts/user/{user_id}", postHandler.GetPostsByUserID).Methods("GET")
	s.HandleFunc("/posts", postHandler.GetAllPosts).Methods("GET")

	// Initialize the server with the configuration settings
	server := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		ReadTimeout:  parseDuration(config.AppConfig.Server.ReadTimeout),
		WriteTimeout: parseDuration(config.AppConfig.Server.WriteTimeout),
		IdleTimeout:  parseDuration(config.AppConfig.Server.IdleTimeout),
		Handler:      r,
	}

	log.Printf("Starting server on port %s", config.AppConfig.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Failed to parse duration: %v", err)
	}
	return d
}
