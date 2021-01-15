package dao

import "github.com/KumKeeHyun/medium-rare/article-service/domain"

type ArticleReplyRepository interface {
	ArticleRepository
	ReplyRepository
	NestedReplyRepository
}

type ArticleRepository interface {
	FindArticleByID(id int) (domain.Article, error)
	FindArticleByPage(page int) ([]domain.Article, error)
	FindArticleByQuery(query string) ([]domain.Article, error)
	FindArticleByIDList(ids []int) ([]domain.Article, error)
	SaveArticle(article domain.Article) (domain.Article, error)
	IncreaseArticleClap(article domain.Article) error
	DeleteArticle(article domain.Article) error
}

type ReplyRepository interface {
	FindReplyByID(id int) (domain.Reply, error)
	FindReplyAll() ([]domain.Reply, error)
	SaveReply(reply domain.Reply) (domain.Reply, error)
	IncreaseReplyClap(reply domain.Reply) error
	DeleteReply(reply domain.Reply) error
}

type NestedReplyRepository interface {
	FindNestedReplyByID(id int) (domain.NestedReply, error)
	FindNestedReplyAll() ([]domain.NestedReply, error)
	SaveNestedReply(reply domain.NestedReply) (domain.NestedReply, error)
	IncreaseNestedReplyClap(reply domain.NestedReply) error
	DeleteNestedReply(reply domain.NestedReply) error
}
