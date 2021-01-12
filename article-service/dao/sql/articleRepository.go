package sql

import (
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type articleRepository struct {
	db *gorm.DB
}

func NewSqlArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{
		db: db,
	}
}

func (ar *articleRepository) FindArticleByID(id int) (result domain.Article, err error) {
	return result, ar.db.Preload("Replies.NestedReplies").Preload(clause.Associations).First(&result, id).Error
}

func (ar *articleRepository) FindArticleByPage(page int) (result []domain.Article, err error) {
	limit := 10
	offset := (limit * (page - 1)) + 1
	return result, ar.db.Limit(limit).Offset(offset).Find(&result).Error
}

func (ar *articleRepository) SearchArticle(query string) (result []domain.Article, err error) {
	return result, ar.db.Where(`match(content) against (? in boolean mode)`, `"`+query+`"`).Find(&result).Error
}

func (ar *articleRepository) SaveArticle(article domain.Article) (domain.Article, error) {
	return article, ar.db.Omit(clause.Associations).Create(&article).Error
}

func (ar *articleRepository) IncreaseArticleClap(article domain.Article) (result domain.Article, err error) {
	return article, ar.db.Model(&article).UpdateColumn("claps", gorm.Expr("claps = claps + ?", 1)).Error
}

func (ar *articleRepository) DeleteArticle(article domain.Article) error {
	return ar.db.Delete(&domain.Article{}, article.ID).Error
}

// func (ar *articleRepository) FindReplyByID(id int) (domain.Reply, error) {

// }

// func (ar *articleRepository) FindReplyAll() ([]domain.Reply, error) {

// }

// func (ar *articleRepository) SaveReply(reply domain.Reply) (domain.Reply, error) {

// }

// func (ar *articleRepository) IncreaseReplyClap(reply domain.Reply) (domain.Reply, error) {

// }

// func (ar *articleRepository) DeleteReply(reply domain.Reply) error {

// }

// func (ar *articleRepository) FindNestedReplyByID(id int) (domain.NestedReply, error) {

// }

// func (ar *articleRepository) FindNestedReplyAll() ([]domain.NestedReply, error) {

// }

// func (ar *articleRepository) SaveNestedReply(reply domain.Reply) (domain.NestedReply, error) {

// }

// func (ar *articleRepository) IncreaseNestedReplyClap(reply domain.Reply) (domain.NestedReply, error) {

// }

// func (ar *articleRepository) DeleteNestedReply(reply domain.NestedReply) error {

// }
