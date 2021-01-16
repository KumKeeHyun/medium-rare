package adapter

import "github.com/KumKeeHyun/medium-rare/user-service/domain"

type CreateUserEvent struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Birth  int    `json:"birth"`
}

func ToAdapterCreateUserEvent(du *domain.User) CreateUserEvent {
	return CreateUserEvent{
		ID:     du.ID,
		Email:  du.Email,
		Name:   du.Name,
		Gender: du.Gender,
		Birth:  du.Birth,
	}
}

type DeleteUserEvent struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func ToAdapterDeleteUserEvent(du *domain.User) DeleteUserEvent {
	return DeleteUserEvent{
		ID:    du.ID,
		Email: du.Email,
		Name:  du.Name,
	}
}
