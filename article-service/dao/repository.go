package dao

import "github.com/KumKeeHyun/medium-rare/article-service/domain"

type ArticleRepository interface {
	FindArticleByID(id int) (domain.Article, error)
	FindArticleByPage(page int) ([]domain.Article, error)
	SearchArticle(query string) ([]domain.Article, error)
	SaveArticle(article domain.Article) (domain.Article, error)
	IncreaseArticleClap(article domain.Article) (domain.Article, error)
	DeleteArticle(article domain.Article) error

	FindReplyByID(id int) (domain.Reply, error)
	FindReplyAll() ([]domain.Reply, error)
	SaveReply(reply domain.Reply) (domain.Reply, error)
	IncreaseReplyClap(reply domain.Reply) (domain.Reply, error)
	DeleteReply(reply domain.Reply) error

	FindNestedReplyByID(id int) (domain.NestedReply, error)
	FindNestedReplyAll() ([]domain.NestedReply, error)
	SaveNestedReply(reply domain.Reply) (domain.NestedReply, error)
	IncreaseNestedReplyClap(reply domain.Reply) (domain.NestedReply, error)
	DeleteNestedReply(reply domain.NestedReply) error
}
