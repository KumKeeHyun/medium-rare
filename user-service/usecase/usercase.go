package usecase

import (
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

type UserUsecase interface {
	FindUsers() ([]domain.User, error)
	FindOne(id int) (domain.User, error)
	Register(user domain.User) (domain.User, error)
	Unregister(user domain.User) error
}

type AuthUsecase interface {
	Login(user domain.User) (domain.TokenPair, error)
	// Logout(id int) error

	// return : new accessToken
	RefreshToken(refreshToken string) (string, error)
}
