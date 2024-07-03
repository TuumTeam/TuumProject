package handlers

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Replace these with your actual credentials
var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "86123716820-pcrs0vpblo97t20qfhprl2tvj5fcammm.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		ClientID:     "Ov23liN4UPJcgqXDCckV",
		ClientSecret: "e7039080350f4210d057aeef67911b76f2023d15",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	facebookOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/facebook/callback",
		ClientID:     "793283989608991",
		ClientSecret: "0c16a3d3373fefa6bc16c9509d0bd135",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
)
