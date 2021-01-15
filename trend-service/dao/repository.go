package dao

import "github.com/KumKeeHyun/medium-rare/trend-service/domain"

type ReadRecordRepository interface {
	Save(record domain.ReadRecord) (domain.ReadRecord, error)
	FindArticlesByQuery(query domain.Query) ([]int, error)
}
