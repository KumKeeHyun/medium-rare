package dao

import "github.com/KumKeeHyun/medium-rare/user-service/domain"

type UserRepository interface {
	FindByID(id int) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByName(name string) (domain.User, error)
	FindAll() ([]domain.User, error)
	Save(user domain.User) (domain.User, error)
	Delete(user domain.User) error
}
