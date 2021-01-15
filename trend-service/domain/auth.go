package domain

import jwt "github.com/dgrijalva/jwt-go"

type AccessClaim struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Birth  int    `json:"birth"`
	jwt.StandardClaims
}
