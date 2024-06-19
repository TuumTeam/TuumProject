package main

import (
	"database/sql"
	"github.com/gorilla/csrf"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"tuum.com/internal/handlers"
	"tuum.com/internal/repositories"
)

func main() {
	// Routes
	http.HandleFunc("/", handlers.RedirectToIndex)
	http.HandleFunc("/home", handlers.RedirectToIndex)
	http.HandleFunc("/tuum", handlers.RedirectToTuum)
	//http.HandleFunc("/login", handlers.RegisterToLogin)
	//http.HandleFunc("/register", handlers.RedirectToRegister)
	http.HandleFunc("/profile", handlers.RedirectToProfile)
	// Autres routes...
	http.Handle("/web/protected/admin.html", handlers.IsAuthenticated(http.HandlerFunc(handlers.ProtectedFileHandler)))

	// Connexion à la base de données SQLite
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Création des repositories et handlers
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Configuration des routes
	http.HandleFunc("/register", userHandler.Signup)
	http.HandleFunc("/login", userHandler.Login)

	// Middleware CSRF protection
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	// Démarrer le serveur
	log.Println("Serveur démarré sur : http://localhost:8080")
	err = http.ListenAndServe(":8080", csrfMiddleware(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}
