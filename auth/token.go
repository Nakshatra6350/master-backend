package auth

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("4ae2e45bd43f808217b41ab969fac8c1")

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]interface{} `json:"custom_claims"`
}

func GetJWT(username string, password string) (map[string]string, error) {
	claims := JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(time.Hour * 6)).Unix(),
		},
		CustomClaims: map[string]interface{}{
			"username": username,
			"password": password,
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenString.SignedString(mySigningKey)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return nil, err
	}

	return map[string]string{
		"access_token": token,
	}, nil
}
