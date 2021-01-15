package domain

import (
	"strconv"
	"time"
)

type Viewed struct {
	UserID    int       `json:"user_id" gorm:"primaryKey;not null"`
	ArticleID int       `json:"article_id" gorm:"primaryKey;not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
}

func ViewsToQuery(vs []Viewed) string {
	if len(vs) == 0 {
		return "0"
	}

	result := strconv.Itoa(vs[0].ArticleID)
	for _, v := range vs[1:] {
		result += "," + strconv.Itoa(v.ArticleID)
	}
	return result
}
