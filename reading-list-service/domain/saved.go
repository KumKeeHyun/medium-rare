package domain

import (
	"strconv"
	"time"
)

type Saved struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;not null"`
	ArticleID int       `json:"article_id" gorm:"primaryKey;not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
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
