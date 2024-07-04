package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
	"tuum.com/internal/auth"
	"tuum.com/internal/database"
)

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
	switch r.URL.Path {
	case "/auth/google/callback":
		oauth2Config = googleOauthConfig
	case "/auth/github/callback":
		oauth2Config = githubOauthConfig
	case "/auth/facebook/callback":
		oauth2Config = facebookOauthConfig
	}

	Otoken, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), Otoken)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userInfo.Body.Close()
	data, _ := ioutil.ReadAll(userInfo.Body)

	var user map[string]interface{}
	if err := json.Unmarshal(data, &user); err != nil {
		http.Error(w, "Failed to unmarshal user info", http.StatusInternalServerError)
		return
	}
	dbUser, err := database.GetUserByEmail(user["email"].(string))
	fmt.Printf("User: %v\n", dbUser)
	if err != nil {
		database.CreateUser(user["name"].(string), user["email"].(string), "")
		dbUser, _ = database.GetUserByEmail(user["email"].(string))
		fmt.Printf("User: %v\n", dbUser)

	}
	err = nil
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
