package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/dgrijalva/jwt-go"
)

func validateUser(user *domain.User) error {
	return nil
}

func hashingPassword(user *domain.User) {
	pw := user.Password + "_kkh"
	res := sha256.Sum256([]byte(pw))
	user.Password = hex.EncodeToString(res[:])
}

func generateTokenPair(user *domain.User) (domain.TokenPair, error) {
	as, err := generateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	rc := domain.RefreshClaim{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    "kkh",
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	rs, err := rt.SignedString([]byte(config.App.JWTSecret))

	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  as,
		RefreshToken: rs,
	}, nil
}

func generateAccessToken(user *domain.User) (string, error) {
	ac := domain.AccessClaim{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Gender: user.Gender,
		Birth:  user.Birth,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			Issuer:    "kkh",
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	as, err := at.SignedString([]byte(config.App.JWTSecret))

	if err != nil {
		return "", err
	}

	return as, nil
}
