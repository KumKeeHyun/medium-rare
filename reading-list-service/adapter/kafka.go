package adapter

import "time"

type UserInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Birth  int    `json:"birth"`
}

type ArticleInfo struct {
	ID         int    `json:"id"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
}

type ReadArticleEvent struct {
	User      UserInfo    `json:"user"`
	Article   ArticleInfo `json:"article"`
	Timestamp time.Time   `json:"timestamp"`
}
