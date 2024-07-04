package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"tuum.com/internal/auth"
	"tuum.com/internal/database"
)

func ExecTmpl(w http.ResponseWriter, tmpl string, data interface{}) {
	err := template.Must(template.ParseFiles(tmpl)).Execute(w, data)
	if err != nil {
		fmt.Printf("Erreur d'execution du template\n")
	}
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/index.html", nil)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		search := r.URL.Query().Get("q")
		db, _ := sql.Open("sqlite3", "./database/forum.db")
		rows, err := db.Query("SELECT name FROM rooms")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		results := []string{}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Fatal(err)
			}
			if strings.Contains(name, search) {
				results = append(results, name)
				// Do something with the matching name
			}
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		response, err := json.Marshal(results)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecTmpl(w, "web/templates/register.html", nil)
	} else {
		if r.FormValue("LogType") == "Login" {
			logBool := database.Login(r.FormValue("email"), r.FormValue("hash"))
			if logBool {
				token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"))
				if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return
				}

				// Set JWT as cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					Expires:  time.Now().Add(24 * time.Hour),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/tuums", http.StatusSeeOther)
			} else {
				http.Error(w, "Login failed", http.StatusUnauthorized)
			}
		} else {
			database.CreateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("hash"))
			token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"))
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}

			// Set JWT as cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})

			http.Redirect(w, r, "/tuums", http.StatusSeeOther)
		}
	}
}

func RedirectToProfile(w http.ResponseWriter, r *http.Request) {
	// Extract session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// If the session cookie is not set, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate JWT
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		// If there is an error in getting the cookie, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get user details from the database
	user := database.GetUserByEmail(claims.Email)
	if err != nil {
		// If there is an error in getting the user, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute the profile template with the user data
	ExecTmpl(w, "web/templates/profile.html", user)
}

func RedirectToTuums(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tuum:")
	if r.Method == http.MethodPost {
		fmt.Println("test")
		if r.FormValue("creationSelector") == "newRoom" {
			fmt.Println("newRoom")
			database.CreateRoom(r.FormValue("title"), r.FormValue("description"))
		} else if r.FormValue("creationSelector") == "newTuum" {
			fmt.Println("newTuum")
			database.CreatePost(r.FormValue("title"), r.FormValue("description"))
		} else {
			fmt.Println("rien")
		}
		fmt.Println("finished")
	}
	ExecTmpl(w, "web/templates/tuum.html", nil)
}

func RedirectToCreate(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/create.html", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Delete the cookie by setting an expired date
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
