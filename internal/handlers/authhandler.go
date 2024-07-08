package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"tuum.com/internal/auth"
	"tuum.com/internal/database"
)

type OauthUser struct {
	Username string
	Email    string
}

func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	var oauthStateString = "pseudo-random"
	url := ""
	switch r.URL.Path {
	case "/auth/google/login":
		url = googleOauthConfig.AuthCodeURL(oauthStateString)
	case "/auth/github/login":
		url = githubOauthConfig.AuthCodeURL(oauthStateString)
	case "/auth/facebook/login":
		url = facebookOauthConfig.AuthCodeURL(oauthStateString)
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	var oauthStateString = "pseudo-random"
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	var oauth2Config *oauth2.Config
	user := make(map[string]interface{})
	oauthUser := OauthUser{}
	switch r.URL.Path {
	case "/auth/google/callback":
		oauth2Config = googleOauthConfig
	case "/auth/facebook/callback":
		oauth2Config = facebookOauthConfig
	case "/auth/github/callback":
		oauth2Config = githubOauthConfig
	}
	Otoken, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}
	switch r.URL.Path {
	case "/auth/google/callback":
		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + Otoken.AccessToken)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return

		}
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &user)
		oauthUser.Username = user["name"].(string)
		oauthUser.Email = user["email"].(string)
	case "/auth/facebook/callback":
		response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + Otoken.AccessToken)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &user)
		oauthUser.Username = user["name"].(string)
		oauthUser.Email = user["email"].(string)
	case "/auth/github/callback":
		oauth2Config = githubOauthConfig
		Otoken, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Code exchange failed", http.StatusInternalServerError)
			return
		}
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", "Bearer "+Otoken.AccessToken)
		response, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &user)
		fmt.Printf("User: %v\n", user)
	}
	dbUser, err := database.GetUserByEmail(user["email"].(string))
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	status, _ := database.GetStatusByEmail(user["email"].(string))
	if "banned" == status {
		w.Write([]byte("<script>alert('User is banned');</script>"))
		ExecTmpl(w, "web/templates/register.html", nil)
		return
	}

	if dbUser.ID == 0 {
		err := database.CreateUser(user["name"].(string), user["email"].(string), "")
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		dbUser, err = database.GetUserByEmail(user["email"].(string))
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateJWT(dbUser.Username, dbUser.Email)
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
