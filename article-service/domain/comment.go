package domain

import "time"

type Reply struct {
	ID            int           `json:"id" gorm:"primaryKey"`
	Claps         int           `json:"claps" gorm:"default:0"`
	Comment       string        `json:"comment" gorm:"not null;type:text"`
	ArticleID     int           `json:"article_id"`
	UserID        int           `json:"user_id" gorm:"not null"`
	UserName      string        `json:"user_name" gorm:"not null"`
	NestedReplies []NestedReply `json:"nested_replies" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type NestedReply struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Claps     int       `json:"claps" gorm:"default:0"`
	Comment   string    `json:"comment" gorm:"not null;type:text"`
	ReplyID   int       `json:"reply_id"`
	UserID    int       `json:"user_id" gorm:"not null"`
	UserName  string    `json:"user_name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
