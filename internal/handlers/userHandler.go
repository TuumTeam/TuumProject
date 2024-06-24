package handlers

import (
	"errors"
	"net/http"
	"time"
	"tuum.com/internal/auth"
	"tuum.com/internal/database"
	"tuum.com/internal/models"
)

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecTmpl(w, r, "web/templates/register.html", nil)
	} else {
		if r.FormValue("LogType") == "Login" {
			logBool, _ := database.Login(r.FormValue("email"), r.FormValue("password"))
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
				ExecTmpl(w, r, "web/templates/register.html", "Invalid email or password")
			}
		} else {
			user := models.User{
				ID:       0,
				Username: r.FormValue("username"),
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			}
			err := database.AddUser(user)
			if err != nil {
				http.Error(w, "Unable to add user to database", http.StatusInternalServerError)
				return
			}

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
	user, err := database.GetUserByEmail(claims.Email)
	if err != nil {
		// If there is an error in getting the user, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Execute the profile template with the user data
	ExecTmpl(w, r, "web/templates/profile.html", user)
}
