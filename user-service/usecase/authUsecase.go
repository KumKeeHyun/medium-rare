package usecase

import (
	"fmt"

	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/dgrijalva/jwt-go"
)

type authUsecase struct {
	ur dao.UserRepository
}

func NewAuthUsecase(ur dao.UserRepository) AuthUsecase {
	return &authUsecase{
		ur: ur,
	}
}

func (au *authUsecase) Login(user domain.User) (domain.TokenPair, error) {
	savedUser, err := au.ur.FindByEmail(user.Email)
	if err != nil {
		return domain.TokenPair{}, err
	}

	hashingPassword(&user)
	if savedUser.Password != user.Password {
		return domain.TokenPair{}, fmt.Errorf("Wrong password for %s", user.Email)
	}

	tokenPair, err := generateTokenPair(&savedUser)
	if err != nil {
		return tokenPair, err
	}

	return tokenPair, nil
}

func (au *authUsecase) RefreshToken(rts string) (string, error) {
	token, err := jwt.ParseWithClaims(rts, &domain.RefreshClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("kkh"), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claim, ok := token.Claims.(*domain.RefreshClaim)
	if !ok {
		return "", fmt.Errorf("Unexpected custom claim type")
	}

	userInfo, err := au.ur.FindByID(claim.ID)
	if err != nil {
		return "", err
	}

	return generateAccessToken(&userInfo)
}
