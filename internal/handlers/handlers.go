package handlers

import (
	"errors"

	"html/template"
	"net/http"
	"time"

	"tuum.com/internal/auth"
	"tuum.com/internal/database"
	"tuum.com/internal/models"
)

func ExecTmpl(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/index.html", nil)
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecTmpl(w, r, "web/templates/register.html", nil)
	} else {
		if r.FormValue("LogType") == "Login" {
			logBool := database.Login(r.FormValue("email"), r.FormValue("password"))
			if logBool {
				// Generate JWT
				token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"), "dark")
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
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Error(w, "Login failed", http.StatusUnauthorized)
			}
		} else {
			database.CreateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
			// Generate JWT
			token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"), "dark")
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

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
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
	ExecTmpl(w, r, "web/templates/profile.html", user)
}

func RedirectToTuums(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/tuums.html", nil)
	if r.Method == http.MethodGet {
		ExecTmpl(w, r, "web/templates/Tuum.html", nil)
	} else {
		if r.FormValue("LogType") == "Login" {
			logBool := database.Login(r.FormValue("email"), r.FormValue("password"))
			if logBool {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Error(w, "Login failed", http.StatusUnauthorized)
			}
		} else {
			if r.FormValue("LogType") == "Login" {
				logBool := database.Login(r.FormValue("email"), r.FormValue("password"))
				if logBool {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				} else {
					http.Error(w, "Login failed", http.StatusUnauthorized)
				}
			} else {
				post := models.Post{
					UserID:    1,
					RoomID:    1,
					Title:     r.FormValue("title"),
					Content:   r.FormValue("content"),
					CreatedAt: time.Now(),
				}
				err := database.AddPost(post)
				if err != nil {
					http.Error(w, "Unable to add user to database", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/tumms", http.StatusSeeOther)
			}
		}
	}
}

func RedirectToCreate(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/create.html", nil)
}
