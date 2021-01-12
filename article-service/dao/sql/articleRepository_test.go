package sql_test

import (
	"testing"

	"github.com/KumKeeHyun/medium-rare/article-service/dao/sql"
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"github.com/KumKeeHyun/medium-rare/article-service/util"
	"github.com/google/go-cmp/cmp"
)

func TestSchema(T *testing.T) {
	_, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}
}

func TestCreateArticle(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	ar := sql.NewSqlArticleRepository(tx)

	result, err := ar.SaveArticle(domain.Article{
		Title:    "hello gorm",
		Content:  "hello my name is kumkeehyun",
		UserID:   1,
		UserName: "test",
	})

	if err != nil {
		T.Error(err)
	}

	saved, err := ar.FindArticleByID(result.ID)
	if err != nil {
		T.Error(err)
	}

	if cmp.Equal(result, saved) {
		T.Fail()
	}
}

func TestSearchArticle(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	ar := sql.NewSqlArticleRepository(db)

	r1, _ := ar.SaveArticle(domain.Article{
		Title:    "hello gorm",
		Content:  "hello my name is kumkeehyun",
		UserID:   1,
		UserName: "test",
	})

	r2, _ := ar.SaveArticle(domain.Article{
		Title:    "hello gorm2",
		Content:  "hello my name is kumkeehyun2",
		UserID:   1,
		UserName: "test",
	})

	result, err := ar.SearchArticle("name")
	if err != nil {
		T.Error(err)
	}

	if len(result) != 2 {
		T.Fail()
	}

	// 같은 트랜젝션 안에선는 전문검색이 안되는 것 같음
	// 그래서 중간에 에러가 안나길 기도하는 메타로 테스팅
	ar.DeleteArticle(r1)
	ar.DeleteArticle(r2)
}
