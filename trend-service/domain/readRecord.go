package domain

import "time"

type ReadRecord struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Birth     int       `json:"birth" gorm:"not null"`
	Gender    string    `json:"gender" gorm:"not null"`
	ArticleID int       `json:"article_id" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
}

type Query struct {
	Gender string
	// 10, 20, 30, 40
	Age int
	// day(1), week(7), month(30), year(365)
	Term int
}

// SQL 쿼리를 위한 값 체크
// 잘못된 값이면 에러보단 기본값("", 0) 으로 수정
func (q *Query) ToValid() {
	if q.Gender != "" && q.Gender != "M" && q.Gender != "F" {
		q.Gender = ""
	}

	switch {
	case q.Age >= 10 && q.Age < 20:
		q.Age = 10
	case q.Age >= 20 && q.Age < 30:
		q.Age = 20
	case q.Age >= 30 && q.Age < 40:
		q.Age = 30
	case q.Age >= 40 && q.Age < 50:
		q.Age = 40
	default:
		q.Age = 0
	}

	switch q.Term {
	case 1, 7, 30, 365:
		// valid
	default:
		q.Term = 0
	}
}
