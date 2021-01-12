package domain

import "time"

type Article struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Title   string `json:"title" gorm:"not null"`
	Content string `json:"content" gorm:"type:text;index:,class:FULLTEXT"`
	Claps   int    `json:"claps" gorm:"default:0"`
	// Content   string    `json:"content" gorm:"type:text;index:,class:FULLTEXT,option:WITH PARSER NGRAM"`
	Replies   []Reply   `json:"replies" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    int       `json:"user_id" gorm:"not null"`
	UserName  string    `json:"user_name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
