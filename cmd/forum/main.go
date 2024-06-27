package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"tuum.com/internal/handlers"
	"tuum.com/pkg/middleware"
)

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/login", handlers.RedirectToLogin)
	r.HandleFunc("/", handlers.RedirectToIndex)
	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.HandleFunc("/logout", handlers.Logout)
	s.HandleFunc("/tuums", handlers.RedirectToTuums)
	s.HandleFunc("/profile", handlers.RedirectToProfile)
	s.HandleFunc("/create", handlers.RedirectToCreate)

	fmt.Println("Server starting at http://localhost:8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
