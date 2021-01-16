package adapter

import (
	"time"

	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birth     int       `json:"birth"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAdapterUser(du *domain.User) User {
	return User{
		ID:        du.ID,
		Email:     du.Email,
		Name:      du.Name,
		Gender:    du.Gender,
		Birth:     du.Birth,
		CreatedAt: du.CreatedAt,
		UpdatedAt: du.UpdatedAt,
	}
}

func ToAdapterUsers(dus []domain.User) []User {
	us := make([]User, 0, len(dus))
	for _, du := range dus {
		us = append(us, ToAdapterUser(&du))
	}
	return us
}
