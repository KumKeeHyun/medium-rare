package controller

import (
	"encoding/json"

	"github.com/KumKeeHyun/medium-rare/trend-service/adapter"
	"github.com/KumKeeHyun/medium-rare/trend-service/dao"
	"github.com/KumKeeHyun/medium-rare/trend-service/domain"
	"go.uber.org/zap"
)

type EventController struct {
	rrr dao.ReadRecordRepository
	log *zap.Logger
}

func NewEventController(rrr dao.ReadRecordRepository, log *zap.Logger) *EventController {
	return &EventController{
		rrr: rrr,
		log: log,
	}
}

func (ec *EventController) ReadArticle(key, value []byte) {
	var rae adapter.ReadArticleEvent
	if err := json.Unmarshal(value, &rae); err != nil {
		ec.log.Error("Fail to unmarshal ReadArticleEvent",
			zap.Error(err))
		return
	}

	rr := domain.ReadRecord{
		Birth:     rae.User.Birth,
		Gender:    rae.User.Gender,
		ArticleID: rae.Article.ID,
		Timestamp: rae.Timestamp,
	}
	if _, err := ec.rrr.Save(rr); err != nil {
		ec.log.Error("Fail to save ReadRecord",
			zap.Error(err))
		return
	}
}
