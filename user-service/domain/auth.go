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
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwibmFtZSI6InRlc3ROYW1lIiwiZ2VuZGVyIjoiTSIsImJpcnRoIjoxOTk5LCJleHAiOjE2MTA3Njk5MTksImlzcyI6ImtraCJ9.6efd8vn9BoFDgBDA9xQCph9xbXvnOaL1DTuYYgEanTQ"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNjEwODEyNTE5LCJpc3MiOiJra2gifQ.9UPSEA_3ngl3HGcch23qXnbO7W-ghfu2Qyqyc01w368"`
}

// AccessToken example for swagger
// not used
type AccessToken struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwibmFtZSI6InRlc3ROYW1lIiwiZ2VuZGVyIjoiTSIsImJpcnRoIjoxOTk5LCJleHAiOjE2MTA3Njk5MTksImlzcyI6ImtraCJ9.6efd8vn9BoFDgBDA9xQCph9xbXvnOaL1DTuYYgEanTQ"`
}

// RefreshToken example for swagger
// not used
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNjEwODEyNTE5LCJpc3MiOiJra2gifQ.9UPSEA_3ngl3HGcch23qXnbO7W-ghfu2Qyqyc01w368"`
}
