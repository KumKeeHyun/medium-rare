package memory

import (
	"fmt"

	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

var (
	store    = map[int]domain.User{}
	sequence = 0
)

func userFilter(store map[int]domain.User, f func(domain.User) bool) []domain.User {
	res := make([]domain.User, 0)
	for _, u := range store {
		if f(u) {
			res = append(res, u)
		}
	}
	return res
}

func userToSlice(store map[int]domain.User) []domain.User {
	res := make([]domain.User, 0, len(store))
	for _, u := range store {
		res = append(res, u)
	}
	return res
}

type memoryUserRepository struct {
}

func NewMemoryUserRepository() dao.UserRepository {
	return &memoryUserRepository{}
}

func (ur *memoryUserRepository) FindByID(id int) (domain.User, error) {
	user, exists := store[id]
	if !exists {
		return domain.User{}, fmt.Errorf("cannot find user where id = %d", id)
	}
	return user, nil
}

func (ur *memoryUserRepository) FindByEmail(email string) (domain.User, error) {
	users := userFilter(store, func(u domain.User) bool {
		return u.Email == email
	})

	if len(users) == 0 {
		return domain.User{}, fmt.Errorf("cannot find user where email = %s", email)
	}
	return users[0], nil
}

func (ur *memoryUserRepository) FindByName(name string) (domain.User, error) {
	users := userFilter(store, func(u domain.User) bool {
		return u.Name == name
	})

	if len(users) == 0 {
		return domain.User{}, fmt.Errorf("cannot find user where name = %s", name)
	}
	return users[0], nil
}

func (ur *memoryUserRepository) FindAll() ([]domain.User, error) {
	return userToSlice(store), nil
}

func (ur *memoryUserRepository) Save(user domain.User) (domain.User, error) {
	sequence++
	user.ID = sequence
	store[user.ID] = user
	return user, nil
}

func (ur *memoryUserRepository) Update(user domain.User) (domain.User, error) {
	u, err := ur.FindByID(user.ID)
	if err != nil {
		return user, err
	}
	nilUser := domain.User{}
	switch {
	case nilUser.Name != user.Name:
		u.Name = user.Name
		fallthrough
	case nilUser.Email != user.Email:
		u.Email = user.Email
		fallthrough
	case nilUser.Password != user.Password:
		u.Password = user.Password
	}
	store[u.ID] = u
	return u, nil
}

func (ur *memoryUserRepository) Delete(user domain.User) error {
	delete(store, user.ID)
	return nil
}
