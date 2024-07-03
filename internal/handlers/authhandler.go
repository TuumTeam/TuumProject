package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

// Redirect to the appropriate OAuth provider's login page
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

// Handle the callback from the OAuth provider
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

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	// Now you can use the token to get the user's information
	client := oauth2Config.Client(context.Background(), token)
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

	// Here you can handle the user information (e.g., create a user in your database)
	fmt.Fprintf(w, "User Info: %+v", user)
}
