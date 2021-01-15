package controller

import (
	"encoding/json"

	"github.com/KumKeeHyun/medium-rare/reading-list-service/adapter"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/dao"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/domain"
	"go.uber.org/zap"
)

type EventController struct {
	rr  dao.ReadingRepository
	log *zap.Logger
}

func NewEventController(rr dao.ReadingRepository, log *zap.Logger) *EventController {
	return &EventController{
		rr:  rr,
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

	v := domain.Viewed{
		UserID:    rae.User.ID,
		ArticleID: rae.Article.ID,
		Timestamp: rae.Timestamp,
	}
	if _, err := ec.rr.SaveViewed(v); err != nil {
		ec.log.Error("Fail to save viewed",
			zap.Any("viewed", v),
			zap.Error(err))
		return
	}
}
