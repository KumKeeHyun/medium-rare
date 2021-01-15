package usecase

import (
	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"go.uber.org/zap"
)

type userUsecase struct {
	ur  dao.UserRepository
	log *zap.Logger
	// producer
}

func NewUserUsecase(ur dao.UserRepository, log *zap.Logger) UserUsecase {
	return &userUsecase{
		ur:  ur,
		log: log,
	}
}

func (uu *userUsecase) FindUsers() ([]domain.User, error) {
	return uu.ur.FindAll()
}

func (uu *userUsecase) FindOne(id int) (domain.User, error) {
	return uu.ur.FindByID(id)
}

func (uu *userUsecase) Register(user domain.User) (domain.User, error) {
	if err := validateUser(&user); err != nil {
		return domain.User{}, err
	}

	hashingPassword(&user)

	return uu.ur.Save(user)
}

func (uu *userUsecase) Unregister(user domain.User) error {
	return uu.ur.Delete(user)
}
