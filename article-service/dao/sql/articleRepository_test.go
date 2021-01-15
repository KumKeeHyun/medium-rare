package sql

import (
	"testing"

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

	ar := NewSqlArticleRepository(tx)

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

	ar := NewSqlArticleRepository(db)

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

	result, err := ar.FindArticleByQuery("name")
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

func TestIncreaseArticleClap(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	tx := db.Begin()
	defer tx.Rollback()
	ar := NewSqlArticleRepository(tx)

	article, _ := ar.SaveArticle(domain.Article{
		Title:    "hello gorm",
		Content:  "hello my name is kumkeehyun",
		UserID:   1,
		UserName: "test",
	})

	err = ar.IncreaseArticleClap(article)
	if err != nil {
		T.Error(err)
	}

	result, _ := ar.FindArticleByID(article.ID)

	if result.Claps != 1 {
		T.Log(result)
		T.Fail()
	}
}

func TestCreateReply(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	tx := db.Begin()
	defer tx.Rollback()
	ar := NewSqlArticleRepository(tx)
	rr := NewSqlReplyRepository(tx)
	nr := NewSqlNestedReplyRepository(tx)

	a, _ := ar.SaveArticle(domain.Article{
		Title:    "hello gorm",
		Content:  "hello my name is kumkeehyun",
		UserID:   1,
		UserName: "test",
	})
	r, _ := rr.SaveReply(domain.Reply{
		Comment:   "haha ok kumkeehyun",
		ArticleID: a.ID,
		UserID:    1,
		UserName:  "test",
	})
	n, err := nr.SaveNestedReply(domain.NestedReply{
		Comment:  "haha ok reply",
		ReplyID:  r.ID,
		UserID:   1,
		UserName: "test",
	})

	if err != nil {
		T.Error(err)
	}

	result, _ := ar.FindArticleByID(a.ID)
	T.Log(result)

	if result.Replies[0].Comment != r.Comment {
		T.Fail()
	}

	if result.Replies[0].NestedReplies[0].Comment != n.Comment {
		T.Fail()
	}
}

func TestFindArticlePage(T *testing.T) {
	db, err := util.BuildMysqlConnection()
	if err != nil {
		T.Error(err)
	}

	tx := db.Begin()
	defer tx.Rollback()
	ar := NewSqlArticleRepository(tx)
	rr := NewSqlReplyRepository(tx)

	a, _ := ar.SaveArticle(domain.Article{
		Title:    "hello gorm",
		Content:  "hello my name is kumkeehyun",
		UserID:   1,
		UserName: "test",
	})
	rr.SaveReply(domain.Reply{
		Comment:   "haha ok kumkeehyun",
		ArticleID: a.ID,
		UserID:    1,
		UserName:  "test",
	})

	result, err := ar.FindArticleByPage(0)
	if err != nil {
		T.Error(err)
	}

	if len(result) != 0 && len(result[0].Replies) != 0 {
		T.Fail()
	}
}
