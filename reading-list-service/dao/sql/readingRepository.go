package sql

import (
	"github.com/KumKeeHyun/medium-rare/reading-list-service/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type readingRepository struct {
	db *gorm.DB
}

func NewSqlReadingRepository(db *gorm.DB) *readingRepository {
	return &readingRepository{
		db: db,
	}
}

func (rr *readingRepository) FindViewedsByUserID(userID int) (result []domain.Viewed, err error) {
	return result, rr.db.Where("user_id = ?", userID).Order("timestamp DESC").Limit(10).Find(&result).Error
}

func (rr *readingRepository) SaveViewed(viewed domain.Viewed) (domain.Viewed, error) {
	return viewed, rr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "article_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"timestamp"}),
	}).Create(&viewed).Error
}

func (rr *readingRepository) FindSavedsByUserID(userID int) (result []domain.Saved, err error) {
	return result, rr.db.Where("user_id = ?", userID).Order("timestamp DESC").Find(&result).Error
}

func (rr *readingRepository) SaveSaved(saved domain.Saved) (domain.Saved, error) {
	return saved, rr.db.Create(&saved).Error
}
