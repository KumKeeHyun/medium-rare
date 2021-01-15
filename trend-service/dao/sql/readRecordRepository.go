package sql

import (
	"time"

	"github.com/KumKeeHyun/medium-rare/trend-service/domain"
	"gorm.io/gorm"
)

type readRecordRepository struct {
	db *gorm.DB
}

func NewSqlReadRecordRepository(db *gorm.DB) *readRecordRepository {
	return &readRecordRepository{
		db: db,
	}
}

func (rrr *readRecordRepository) Save(record domain.ReadRecord) (result domain.ReadRecord, err error) {
	return result, rrr.db.Create(&record).Error
}

func (rrr *readRecordRepository) FindArticlesByQuery(query domain.Query) (ids []int, err error) {
	query.ToValid()
	sqlQuery := rrr.db.Model(&domain.ReadRecord{}).Select("article_id")

	// 성별, 나이, 읽은 시간에 따른 필터링
	if query.Gender != "" {
		sqlQuery = sqlQuery.Where("gender = ?", query.Gender)
	}
	if query.Age != 0 {
		to := time.Now().Year() - query.Age + 1
		from := to - 9
		sqlQuery = sqlQuery.Where("birth BETWEEN ? AND ?", from, to)
	}
	if query.Term != 0 {
		since := time.Now().AddDate(0, 0, query.Term)
		sqlQuery = sqlQuery.Where("timestamp > ?", since)
	}

	return ids, sqlQuery.Group("article_id").Limit(10).Order("count(*) DESC").Scan(&ids).Error
}
