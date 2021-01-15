package sql

import (
	"testing"
	"time"

	"github.com/KumKeeHyun/medium-rare/reading-list-service/domain"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/util"
)

func TestSaveViewed(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	rr := NewSqlReadingRepository(tx)

	v := domain.Viewed{
		UserID:    1,
		ArticleID: 59,
		Timestamp: time.Now(),
	}

	rr.SaveViewed(v)

	v.Timestamp = time.Now().Add(time.Second * 1)
	rr.SaveViewed(v)

	result, _ := rr.FindViewedsByUserID(1)
	if len(result) != 1 {
		T.Log(result)
		T.Fail()
	}
}
