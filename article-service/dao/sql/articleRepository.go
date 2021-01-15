package sql

import (
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type articleReplyRepository struct {
	articleRepository
	replyRepository
	nestedReplyRepository
}

func NewSqlArticleReplyRepository(db *gorm.DB) *articleReplyRepository {
	return &articleReplyRepository{
		articleRepository{db},
		replyRepository{db},
		nestedReplyRepository{db},
	}
}

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

func (ar *articleRepository) FindArticleByQuery(query string) (result []domain.Article, err error) {
	return result, ar.db.Where(`match(content) against (? in boolean mode)`, `"`+query+`"`).Find(&result).Error
}

func (ar *articleRepository) FindArticleByIDList(ids []int) (result []domain.Article, err error) {
	return result, ar.db.Where("id IN ?", ids).Find(&result).Error
}

func (ar *articleRepository) SaveArticle(article domain.Article) (domain.Article, error) {
	return article, ar.db.Omit(clause.Associations).Create(&article).Error
}

func (ar *articleRepository) IncreaseArticleClap(article domain.Article) error {
	return ar.db.Model(&article).UpdateColumn("claps", gorm.Expr("claps + ?", 1)).Error
}

func (ar *articleRepository) DeleteArticle(article domain.Article) error {
	return ar.db.Delete(&domain.Article{}, article.ID).Error
}
