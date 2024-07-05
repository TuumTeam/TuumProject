package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateSessionToken() string {
	var token string
	rand.Seed(time.Now().UnixNano())
	minValue := 33
	maxValue := 126
	fmt.Println()
	for i := 0; i < 16; i++ {
		token += string(rune(rand.Intn(maxValue-minValue+1) + minValue))
	}
	return token
}

var jwtKey = []byte(generateSessionToken())

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(username, email string) (string, error) {
	if jwtKey == nil {
		jwtKey = []byte(generateSessionToken())
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	returnValue, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return "", err
	}
	return returnValue, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
