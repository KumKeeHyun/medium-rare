package sql

import (
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type nestedReplyRepository struct {
	db *gorm.DB
}

func NewSqlNestedReplyRepository(db *gorm.DB) *nestedReplyRepository {
	return &nestedReplyRepository{
		db: db,
	}
}

func (nr *nestedReplyRepository) FindNestedReplyByID(id int) (result domain.NestedReply, err error) {
	return result, nr.db.Preload(clause.Associations).First(&result, id).Error
}

func (nr *nestedReplyRepository) FindNestedReplyAll() (result []domain.NestedReply, err error) {
	return result, nr.db.Find(&result).Error
}

func (nr *nestedReplyRepository) SaveNestedReply(reply domain.NestedReply) (domain.NestedReply, error) {
	return reply, nr.db.Create(&reply).Error
}

func (nr *nestedReplyRepository) IncreaseNestedReplyClap(reply domain.NestedReply) error {
	return nr.db.Model(&reply).UpdateColumn("claps", gorm.Expr("claps + ?", 1)).Error

}

func (nr *nestedReplyRepository) DeleteNestedReply(reply domain.NestedReply) error {
	return nr.db.Delete(&domain.NestedReply{}, reply.ID).Error
}
