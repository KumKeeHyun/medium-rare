package sql

import (
	"testing"
	"time"

	"github.com/KumKeeHyun/medium-rare/trend-service/domain"
	"github.com/KumKeeHyun/medium-rare/trend-service/util"
)

func TestSave(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	rrr := NewSqlReadRecordRepository(db)

	r := domain.ReadRecord{
		Birth:     1999,
		Gender:    "M",
		ArticleID: 1,
		Timestamp: time.Now(),
	}
	for i := 1; i < 13; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}
	for i := 6; i < 13; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}
	for i := 2; i < 7; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}

	r.Birth = 1989
	for i := 1; i < 4; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}
	for i := 7; i < 9; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}

	r.Birth = 2009
	for i := 1; i < 7; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}
	for i := 2; i < 7; i++ {
		r.ArticleID = i
		rrr.Save(r)
	}

}

func TestFindArticlesByQuery(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	rrr := NewSqlReadRecordRepository(db)

	query := domain.Query{
		Gender: "M",
		Age:    20,
	}

	result, _ := rrr.FindArticlesByQuery(query)
	T.Log(result)

	// to show result
	T.Fail()
}
