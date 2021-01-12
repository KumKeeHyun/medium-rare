package sql

import (
	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"gorm.io/gorm"
)

type sqlUserRepository struct {
	db *gorm.DB
}

func NewSqlUserRepository(db *gorm.DB) dao.UserRepository {
	return &sqlUserRepository{
		db: db,
	}
}

func (ur *sqlUserRepository) FindByID(id int) (result domain.User, err error) {
	return result, ur.db.First(&result, id).Error
}

func (ur *sqlUserRepository) FindByEmail(email string) (result domain.User, err error) {
	return result, ur.db.Where("email = ?", email).Find(&result).Error
}

func (ur *sqlUserRepository) FindByName(name string) (result domain.User, err error) {
	return result, ur.db.Where("name = ?", name).Find(&result).Error
}

func (ur *sqlUserRepository) FindAll() (result []domain.User, err error) {
	return result, ur.db.Find(&result).Error
}

func (ur *sqlUserRepository) Save(user domain.User) (domain.User, error) {
	return user, ur.db.Create(&user).Error
}

func (ur *sqlUserRepository) Update(user domain.User) (domain.User, error) {
	return user, ur.db.Model(&user).Updates(user).Error
}

func (ur *sqlUserRepository) Delete(user domain.User) error {
	return ur.db.Delete(&domain.User{}, user.ID).Error
}
