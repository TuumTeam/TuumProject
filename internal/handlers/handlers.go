package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"tuum.com/internal/auth"
	"tuum.com/internal/database"
)

func ExecTmpl(w http.ResponseWriter, tmplPath string, data interface{}) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err // Retourne l'erreur si le fichier ne peut pas être parsé
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return err // Retourne l'erreur si l'exécution échoue
	}
	return nil
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
		if "banned" == database.GetStatusByEmail(r.FormValue("email")) {
			w.Write([]byte("<script>alert('User is banned');</script>"))
			ExecTmpl(w, "web/templates/register.html", nil)
			return
		}
		if r.FormValue("LogType") == "Login" {
			logBool, _ := database.Login(r.FormValue("email"), r.FormValue("hash"))
			if logBool {
				token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"))
				if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return
				}

				// Set JWT as cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",                // Cookie name
					Value:    token,                          // JWT token as the value
					Path:     "/",                            // Set cookie for entire website
					Expires:  time.Now().Add(24 * time.Hour), // Set expiration time
					HttpOnly: true,                           // Make cookie inaccessible to JavaScript
					Secure:   true,                           // Set to true if serving over HTTPS
				})
				http.Redirect(w, r, "/tuums", http.StatusSeeOther)
			} else {
				http.Error(w, "Login failed", http.StatusUnauthorized)
			}
		} else {
			if database.CheckUserExists(r.FormValue("username"), r.FormValue("email")) {
				//w.Write([]byte("<script>document.getElementById('error_message').innerText = 'Name or Email already exists';</script>"))
				w.Write([]byte("<script>alert('Name or Email already exists');</script>"))
				ExecTmpl(w, "web/templates/register.html", nil)
			} else {

				database.CreateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("hash"))
				token, err := auth.GenerateJWT(r.FormValue("username"), r.FormValue("email"))
				if err != nil {
					http.Error(w, "Failed to generate token", http.StatusInternalServerError)
					return
				}

				// Set JWT as cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",                // Cookie name
					Value:    token,                          // JWT token as the value
					Path:     "/",                            // Set cookie for entire website
					Expires:  time.Now().Add(24 * time.Hour), // Set expiration time
					HttpOnly: true,                           // Make cookie inaccessible to JavaScript
					Secure:   true,                           // Set to true if serving over HTTPS
				})

				http.Redirect(w, r, "/tuums", http.StatusSeeOther)
			}
		}
	}
}

func RedirectToAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Handle error, redirect to login if no cookie
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate JWT
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		// Handle error, invalid token
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user details from the database
	user, err := database.GetUserByEmail(claims.Email)
	if err != nil {
		// Handle error, user not found
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	fmt.Println(user.Status)
	// Check if the user is authorized to access the admin page
	if user.Status != "admin" {
		// Redirect to profile or another appropriate page if not authorized
		http.Redirect(w, r, "/profile", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		searchType := r.FormValue("searchType")
		if searchType == "user" {
			users := database.GetUsers()
			ExecTmpl(w, "web/templates/admin.html", users)
		} else if searchType == "room" {
			rooms := database.GetRooms()
			ExecTmpl(w, "web/templates/admin.html", rooms)
		} else if searchType == "post" {
			posts := database.GetPosts()
			ExecTmpl(w, "web/templates/admin.html", posts)
		} else {
			ExecTmpl(w, "web/templates/admin.html", nil)
		}
	} else {
		users := database.GetUsers()
		ExecTmpl(w, "web/templates/admin.html", users)
	}

}
func BanProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//searchType := r.FormValue("searchType")
		//if searchType == "user" {
		fmt.Println("user test")
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		IdBanished := r.Form["IdBanished"]
		fmt.Println("IdBanished:", IdBanished)
		for _, id := range IdBanished {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(idInt)
			database.ChangeStatusUserByemail("banned", database.GetUser(idInt).Email)

		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		w.Write([]byte("<script>alert('User banned');</script>"))
		/*} else if searchType == "room" {
			IdDelete := r.FormValue("IdDelete")
			for i := 0; i < len(IdDelete); i++ {
				id, err := strconv.Atoi(string(IdDelete[i]))
				if err != nil {
					fmt.Println(err)
				}
				database.DeleteRoom(id)
			}
		} else if searchType == "post" {
			IdDelete := r.FormValue("IdDelete")
			for i := 0; i < len(IdDelete); i++ {
				id, err := strconv.Atoi(string(IdDelete[i]))
				if err != nil {
					fmt.Println(err)
				}
				database.DeletePost(id)
			}
		}*/
		//http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	user, _ := database.GetUserByEmail(claims.Email)

	// Execute the profile template with the user data
	ExecTmpl(w, "web/templates/profile.html", user)
}

func RedirectToTuums(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.FormValue("creationSelector") == "newRoom" {
			if database.CheckRoomExists(r.FormValue("title")) {
				w.Write([]byte("<script>alert('the Room already exists');</script>"))
			} else {
				database.CreateRoom(r.FormValue("title"), r.FormValue("description"))
			}
		} else if r.FormValue("creationSelector") == "newTuum" {
			token := r.Header.Get("Authorization")
			claims, err := auth.ValidateJWT(token)
			if err != nil {
				// Log l'erreur et envoie une réponse d'erreur
				log.Printf("Erreur de validation JWT : %v", err)
				http.Error(w, "Non autorisé", http.StatusUnauthorized)
				return
			}
			userEmail := claims.Email
			User, err := database.GetUserByEmail(userEmail)
			if err != nil {
				fmt.Println(err)
			}
			idUser := User.ID
			nameRoom := r.FormValue("searchRoom")
			fmt.Println("name:", nameRoom)
			idRoom := database.GetRoomIdByName(nameRoom)
			fmt.Println(idRoom)
			database.CreatePost(idUser, idRoom, r.FormValue("title"), r.FormValue("description"))
		} else {
			fmt.Println("rien")
		}
		fmt.Println("finished")
	}
	dataBase := database.GetDatabaseForTuum()
	ExecTmpl(w, "web/templates/tuums.html", dataBase)
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
	http.Redirect(w, r, "/tuums", http.StatusSeeOther)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		// Call the function to delete the user account
		err := database.DeleteAccountByEmail(email)
		if err != nil {
			http.Error(w, "Failed to delete account", http.StatusInternalServerError)
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		// Redirect to logout or confirmation page after deletion
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate JWT and extract claims
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to validate session", http.StatusInternalServerError)
		return
	}

	// Fetch user details from the database
	user, err := database.GetUserByEmail(claims.Email)
	if err != nil {
		log.Printf("Failed to fetch user details: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the profile template with the user data
	tmpl, err := template.ParseFiles("web/templates/profile.html")
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func RedirectToAdmin(w http.ResponseWriter, r *http.Request) {
	// Ensure this handler only responds to POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Handle error, redirect to login if no cookie
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate JWT
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		// Handle error, invalid token
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user details from the database
	user, err := database.GetUserByEmail(claims.Email)
	if err != nil {
		// Handle error, user not found
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// Check if the user is authorized to access the admin page
	if user.Status != "admin" {
		// Redirect to profile or another appropriate page if not authorized
		http.Redirect(w, r, "/profile", http.StatusForbidden)
		return
	}

	// Execute the admin template with the user data
	ExecTmpl(w, "web/templates/admin.html", user)
}
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Check Request Method
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 2: Authenticate User
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Step 3: Authorize User
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := database.GetUserByEmail(claims.Email)
	if err != nil || user.Status != "admin" {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	// Step 4: Fetch Data (if needed)
	// Example: data := fetchDataForAdmin()

	// Step 5: Render Template
	ExecTmpl(w, "web/templates/admin.html", user)
}
