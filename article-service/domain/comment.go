package domain

import "time"

type Reply struct {
	ID            int           `json:"id" gorm:"primaryKey" example:"1"`
	Claps         int           `json:"claps" gorm:"default:0" example:"5"`
	Comment       string        `json:"comment" gorm:"not null;type:text" example:"example comment..."`
	ArticleID     int           `json:"article_id" example:"1"`
	UserID        int           `json:"user_id" gorm:"not null" example:"1"`
	UserName      string        `json:"user_name" gorm:"not null" example:"test"`
	NestedReplies []NestedReply `json:"nested_replies" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt     time.Time     `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt     time.Time     `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

type NestedReply struct {
	ID        int       `json:"id" gorm:"primaryKey" example:"1"`
	Claps     int       `json:"claps" gorm:"default:0" example:"2"`
	Comment   string    `json:"comment" gorm:"not null;type:text" example:"example nested comment..."`
	ReplyID   int       `json:"reply_id" example:"1"`
	UserID    int       `json:"user_id" gorm:"not null" example:"1"`
	UserName  string    `json:"user_name" gorm:"not null" example:"test"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

// ReplyNoNested example for swagger
// not used
type ReplyNoNested struct {
	ID        int       `json:"id" example:"1"`
	Claps     int       `json:"claps" example:"5"`
	Comment   string    `json:"comment" example:"example comment..."`
	ArticleID int       `json:"article_id" example:"1"`
	UserID    int       `json:"user_id" example:"1"`
	UserName  string    `json:"user_name" example:"test"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

// CreateReply example for swagger
// not used
type CreateReply struct {
	Comment string `json:"comment" example:"example comment..."`
}
