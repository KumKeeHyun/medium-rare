package sql

import (
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type replyRepository struct {
	db *gorm.DB
}

func NewSqlReplyRepository(db *gorm.DB) *replyRepository {
	return &replyRepository{
		db: db,
	}
}

func (rr *replyRepository) FindReplyByID(id int) (result domain.Reply, err error) {
	return result, rr.db.Preload(clause.Associations).First(&result, id).Error
}
func (rr *replyRepository) FindReplyAll() (result []domain.Reply, err error) {
	return result, rr.db.Find(&result).Error

}
func (rr *replyRepository) SaveReply(reply domain.Reply) (domain.Reply, error) {
	return reply, rr.db.Omit(clause.Associations).Create(&reply).Error
}
func (rr *replyRepository) IncreaseReplyClap(reply domain.Reply) error {
	return rr.db.Model(&reply).UpdateColumn("claps", gorm.Expr("claps + ?", 1)).Error
}
func (rr *replyRepository) DeleteReply(reply domain.Reply) error {
	return rr.db.Delete(&domain.Reply{}, reply.ID).Error
}
