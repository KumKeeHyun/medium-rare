package domain

import "time"

type Article struct {
	ID        int       `json:"id" gorm:"primaryKey" example:"1"`
	Title     string    `json:"title" gorm:"not null" example:"example title"`
	Content   string    `json:"content" gorm:"type:text;index:,class:FULLTEXT" example:"example content..."`
	Claps     int       `json:"claps" gorm:"default:0" example:"123"`
	Replies   []Reply   `json:"replies" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    int       `json:"user_id" gorm:"not null" example:"1"`
	UserName  string    `json:"user_name" gorm:"not null" example:"test"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

// ArticleList example for swagger
// not used
type ArticleList struct {
	ArticleList []ArticleNoReply `json:"article_list"`
}

// ArticleNoReply example for swagger
// not used
type ArticleNoReply struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"example title"`
	Content   string    `json:"content" example:"example content..."`
	Claps     int       `json:"claps" example:"123"`
	UserID    int       `json:"user_id" example:"1"`
	UserName  string    `json:"user_name" example:"test"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-15T09:44:35.151+09:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-01-15T09:44:35.151+09:00"`
}

// ArticleForSingle example for swagger
// not used
type ArticleForSingle struct {
	ArticleList Article `json:"article"`
}

// CreateArticle example for swagger
// not used
type CreateArticle struct {
	Title   string `json:"title" example:"example title"`
	Content string `json:"content" example:"example content..."`
}
