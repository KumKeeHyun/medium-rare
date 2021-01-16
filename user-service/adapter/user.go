package adapter

import (
	"time"

	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

type User struct {
	ID        int       `json:"id" example:"1"`
	Email     string    `json:"email" example:"test@example.com"`
	Name      string    `json:"name" example:"test"`
	Gender    string    `json:"gender" example:"M"`
	Birth     int       `json:"birth" example:"1999"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
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
