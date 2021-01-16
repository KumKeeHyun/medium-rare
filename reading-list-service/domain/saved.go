package domain

import (
	"strconv"
	"time"
)

type Saved struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;not null" example:"1"`
	ArticleID int       `json:"article_id" gorm:"primaryKey;not null" example:"1"`
	Timestamp time.Time `json:"timestamp" gorm:"not null" example:"2021-01-15T09:44:35.151+09:00"`
}

func SavedsToQuery(ss []Saved) string {
	if len(ss) == 0 {
		return "0"
	}

	result := strconv.Itoa(ss[0].ArticleID)
	for _, s := range ss[1:] {
		result += "," + strconv.Itoa(s.ArticleID)
	}
	return result
}
