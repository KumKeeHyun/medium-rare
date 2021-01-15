package domain

import "github.com/dgrijalva/jwt-go"

// AccessClaim ...
// 다른 마이크로서비스에서도 같은 형식으로 정의해서 써야함
type AccessClaim struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Birth  int    `json:"birth"`
	jwt.StandardClaims
}

type RefreshClaim struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
