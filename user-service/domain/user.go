package domain

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey" example:"1"`
	Email     string    `json:"email" gorm:"unique;not null" example:"test@example.com"`
	Password  string    `json:"password" gorm:"not null" example:"testpw"`
	Name      string    `json:"name" gorm:"not null" example:"test"`
	Gender    string    `json:"gender" gorm:"not null" example:"M"`
	Birth     int       `json:"birth" gorm:"not null" example:"1999"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

// LoginUser example for swagger
// not used
type LoginUser struct {
	Email    string `json:"email" gorm:"unique;not null" example:"test@example.com"`
	Password string `json:"password" gorm:"not null" example:"testpw"`
}

// CreateUser example for swagger
// not used
type CreateUser struct {
	Email    string `json:"email" gorm:"unique;not null" example:"test@example.com"`
	Password string `json:"password" gorm:"not null" example:"testpw"`
	Name     string `json:"name" gorm:"not null" example:"test"`
	Gender   string `json:"gender" gorm:"not null" example:"M"`
	Birth    int    `json:"birth" gorm:"not null" example:"1999"`
}
