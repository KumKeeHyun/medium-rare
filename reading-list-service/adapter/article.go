package adapter

type Article struct {
	ID       int    `json:"id" example:"1"`
	Title    string `json:"title" example:"example title"`
	Content  string `json:"content" example:"example contents"`
	UserID   int    `json:"user_id" example:"1"`
	UserName string `json:"user_name" example:"test"`
}

// ArticleList example for swagger
// not used
type ArticleList struct {
	ArticleList []Article `json:"article_list"`
}
