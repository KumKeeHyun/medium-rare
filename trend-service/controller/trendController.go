package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KumKeeHyun/medium-rare/trend-service/adapter"
	"github.com/KumKeeHyun/medium-rare/trend-service/config"
	"github.com/KumKeeHyun/medium-rare/trend-service/dao"
	"github.com/KumKeeHyun/medium-rare/trend-service/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type TrendController struct {
	rrr dao.ReadRecordRepository
	log *zap.Logger
}

func NewTrendController(rrr dao.ReadRecordRepository, log *zap.Logger) *TrendController {
	return &TrendController{
		rrr: rrr,
		log: log,
	}
}

func (tc *TrendController) ListTrend(c *gin.Context) {
	q := domain.Query{}
	trend, err := tc.rrr.FindArticlesByQuery(q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	var articles []adapter.Article
	ids := func(intIDs []int) string {
		if len(intIDs) == 0 {
			return "0"
		}

		result := strconv.Itoa(intIDs[0])
		for _, id := range intIDs[1:] {
			result += "," + strconv.Itoa(id)
		}
		return result
	}(trend)
	url := fmt.Sprintf("http://%s%s", config.App.ArticleConfig.Address, config.App.ArticleConfig.URL)

	resp, err := resty.New().R().SetHeader("Content-Type", "application/json").
		SetQueryParam("ids", ids).SetResult(&articles).
		Get(url)

	if err != nil || resp.StatusCode() != 200 {
		tc.log.Error("Fail to get article list",
			zap.String("endpoint", url),
			zap.Any("query", ids),
			zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Fail to get article list from article-service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trend": articles})
}

func (tc *TrendController) ListTrendForUser(c *gin.Context) {
	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	q := domain.Query{
		Gender: claims.Gender,
		Age:    time.Now().Year() - claims.Birth + 1,
		Term:   7,
	}
	trend, err := tc.rrr.FindArticlesByQuery(q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	var articles []adapter.Article
	ids := func(intIDs []int) string {
		if len(intIDs) == 0 {
			return "0"
		}

		result := strconv.Itoa(intIDs[0])
		for _, id := range intIDs[1:] {
			result += "," + strconv.Itoa(id)
		}
		return result
	}(trend)
	url := fmt.Sprintf("http://%s%s", config.App.ArticleConfig.Address, config.App.ArticleConfig.URL)

	resp, err := resty.New().R().SetHeader("Content-Type", "application/json").
		SetQueryParam("ids", ids).SetResult(&articles).
		Get(url)

	if err != nil || resp.StatusCode() != 200 {
		tc.log.Error("Fail to get article list",
			zap.String("endpoint", url),
			zap.Any("query", ids),
			zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Fail to get article list from article-service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trend": articles})
}

func getAccessClaims(c *gin.Context) (*domain.AccessClaim, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, exists
	}
	accessClaims := claims.(*domain.AccessClaim)
	return accessClaims, exists
}
