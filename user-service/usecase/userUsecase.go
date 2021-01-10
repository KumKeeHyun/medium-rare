package usecase

import (
	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

type userUsecase struct {
	ur dao.UserRepository
	// producer
}

func NewUserUsecase(ur dao.UserRepository) UserUsecase {
	return &userUsecase{
		ur: ur,
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
	// 이메일 인증 전(init), 후(general)
	// domain에 ENUM 처럼 정의해야함
	user.Role = "init"

	return uu.ur.Save(user)
}

func (uu *userUsecase) Unregister(user domain.User) error {
	return uu.ur.Delete(user)
}
