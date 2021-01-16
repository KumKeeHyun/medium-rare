package adapter

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

// ArticleList example for swagger
// not used
type ArticleList struct {
	ArticleList []Article `json:"article_list"`
}
