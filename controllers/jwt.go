package controllers

import (
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

type Jwt struct {
	TokenAuth *jwtauth.JWTAuth
}

var jwt Jwt

func NewJwt() *Jwt {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("\033[31mJWT_SECRET key not set in environment variable\033[0m")
		secret = "secret"
	}

	jwt.TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	return &jwt
}

func GetTokenAuth() *Jwt {
	return &jwt
}

func (j *Jwt) GetToken(claims map[string]interface{}) string {
	_, tokenString, _ := j.TokenAuth.Encode(claims)
	return tokenString
}
