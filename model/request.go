package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}
