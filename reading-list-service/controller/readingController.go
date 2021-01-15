package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KumKeeHyun/medium-rare/reading-list-service/adapter"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/config"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/dao"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type ReadingController struct {
	rr  dao.ReadingRepository
	log *zap.Logger
}

func NewReadingController(rr dao.ReadingRepository, log *zap.Logger) *ReadingController {
	return &ReadingController{
		rr:  rr,
		log: log,
	}
}

func (rc *ReadingController) ListRecent(c *gin.Context) {
	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	vieweds, err := rc.rr.FindViewedsByUserID(claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	var articles []adapter.Article
	ids := domain.ViewsToQuery(vieweds)
	url := fmt.Sprintf("http://%s%s", config.App.ArticleConfig.Address, config.App.ArticleConfig.URL)

	resp, err := resty.New().R().SetHeader("Content-Type", "application/json").
		SetQueryParam("ids", ids).SetResult(&articles).
		Get(url)

	if err != nil || resp.StatusCode() != 200 {
		rc.log.Error("Fail to get article list",
			zap.String("endpoint", url),
			zap.Any("query", ids),
			zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Fail to get article list from article-service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recent": articles})
}

func (rc *ReadingController) SaveArticle(c *gin.Context) {
	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	strID, exists := c.GetPostForm("article_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "require postform 'aritlce_id'"})
		return
	}
	articleID, err := strconv.Atoi(strID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "'aritlce_id' require int"})
		return
	}

	saved := domain.Saved{claims.ID, articleID, time.Now()}
	result, err := rc.rr.SaveSaved(saved)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (rc *ReadingController) ListSaved(c *gin.Context) {
	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	saveds, err := rc.rr.FindSavedsByUserID(claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	var articles []adapter.Article
	ids := domain.SavedsToQuery(saveds)
	url := fmt.Sprintf("http://%s%s", config.App.ArticleConfig.Address, config.App.ArticleConfig.URL)

	resp, err := resty.New().R().SetHeader("Content-Type", "application/json").
		SetQueryParam("ids", ids).SetResult(&articles).
		Get(url)

	if err != nil || resp.StatusCode() != 200 {
		rc.log.Error("Fail to get article list",
			zap.String("endpoint", url),
			zap.Any("query", ids),
			zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": "Fail to get article list from article-service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"saved": articles})
}

func getAccessClaims(c *gin.Context) (*domain.AccessClaim, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, exists
	}
	accessClaims := claims.(*domain.AccessClaim)
	return accessClaims, exists
}
