package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KumKeeHyun/medium-rare/article-service/adapter"
	"github.com/KumKeeHyun/medium-rare/article-service/dao"
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleController struct {
	arr      dao.ArticleReplyRepository
	producer sarama.SyncProducer
	log      *zap.Logger
}

func NewArticleController(arr dao.ArticleReplyRepository,
	sp sarama.SyncProducer,
	log *zap.Logger) *ArticleController {
	return &ArticleController{
		arr:      arr,
		producer: sp,
		log:      log,
	}
}

// ListArticles swagger
// @Summary list all articles
// @Accept json
// @Produce json
// @Param p query string false "page num"
// @Success 200 {object} domain.ArticleList
// @Failure 500 {object} controller.HttpError
// @Router /v1/articles [get]
func (ac *ArticleController) ListArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		page = 0
	}

	articles, err := ac.arr.FindArticleByPage(page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article_list": articles})
}

// ListArticlesByIDList swagger
// @Summary list articles where article id is in ids
// @Accept json
// @Produce json
// @Param ids query string true "article id list separated by comma" Enums("1,2,3")
// @Success 200 {object} []domain.ArticleNoReply
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /v1/articles/list [get]
func (ac *ArticleController) ListArticlesByIDList(c *gin.Context) {
	strIDs := strings.Split(c.Query("ids"), ",")
	if len(strIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "must request with QueryArray named 'ids'"})
		return
	}

	ids := make([]int, 0, len(strIDs))
	for _, strID := range strIDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "'ids' QueryArray require int array"})
			return
		}
		ids = append(ids, id)
	}

	articles, err := ac.arr.FindArticleByIDList(ids)
	if err != nil {
		ac.log.Error("Fail to get article list from db",
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// SearchArticles swagger
// @Summary list articles search by q
// @Accept json
// @Produce json
// @Param q query string true "some word in article content"
// @Success 200 {object} []domain.ArticleList
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Router /v1/articles/search [get]
func (ac *ArticleController) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "must request with Query named 'q'"})
		return
	}

	articles, err := ac.arr.FindArticleByQuery(query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article_list": articles})
}

// GetArticle swagger
// @Summary show article
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Success 200 {object} domain.ArticleForSingle
// @Failure 400 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id} [get]
func (ac *ArticleController) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("article-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	article, err := ac.arr.FindArticleByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	if claims, exists := getAccessClaims(c); exists {
		ac.log.Info("raise read-article event",
			zap.Int("user_id", claims.ID),
			zap.Int("article_id", id))

		rae := adapter.ReadArticleEvent{
			User: adapter.UserInfo{
				ID:     claims.ID,
				Name:   claims.Name,
				Gender: claims.Gender,
				Birth:  claims.Birth,
			},
			Article: adapter.ArticleInfo{
				ID:         id,
				AuthorID:   article.UserID,
				AuthorName: article.UserName,
			},
			Timestamp: time.Now(),
		}

		mshld, _ := json.Marshal(rae)
		ac.producer.SendMessage(&sarama.ProducerMessage{
			Topic: "read-article",
			Value: sarama.ByteEncoder(mshld),
		})
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// CreateArticle swagger
// @Summary create article
// @Accept json
// @Produce json
// @Param article body domain.CreateArticle true "title and content"
// @Success 200 {object} domain.ArticleNoReply
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article [post]
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	var article domain.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}
	article.UserID = claims.ID
	article.UserName = claims.Name

	articleResult, err := ac.arr.SaveArticle(article)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articleResult)
}

// DeleteArticle swagger
// @Summary delete article
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Success 200 {object} integer "article id"
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 403 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id} [delete]
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	article, err := ac.arr.FindArticleByID(articleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	if article.UserID != claims.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "this article is not yours"})
		return
	}

	if err := ac.arr.DeleteArticle(article); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	// TODO : raise delete-article event

	c.JSON(http.StatusOK, articleID)
}

// CreateReply swagger
// @Summary create rely
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Param article body domain.CreateReply true "comment"
// @Success 200 {object} domain.ReplyNoNested
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id}/reply [post]
func (ac *ArticleController) CreateReply(c *gin.Context) {
	var reply domain.Reply

	articleID, err := strconv.Atoi(c.Param("article-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&reply); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	reply.ArticleID = articleID
	reply.UserID = claims.ID
	reply.UserName = claims.Name

	replyResult, err := ac.arr.SaveReply(reply)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replyResult)
}

// DeleteReply swagger
// @Summary delete reply
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Param reply-id path int true "reply's id"
// @Success 200 {object} integer "reply id"
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 403 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id}/reply/{reply-id} [delete]
func (ac *ArticleController) DeleteReply(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("article-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	replyID, err := strconv.Atoi(c.Param("reply-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	reply, err := ac.arr.FindReplyByID(replyID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	if reply.ArticleID != articleID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "reply's articleID is not mathed with param"})
		return
	}

	if reply.UserID != claims.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "this reply is not yours"})
		return
	}

	if err := ac.arr.DeleteReply(reply); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replyID)
}

// CreateNestedReply swagger
// @Summary create nested rely
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Param reply-id path int true "reply's id"
// @Param reply body domain.CreateReply true "comment"
// @Success 200 {object} domain.NestedReply
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id}/reply/{reply-id}/nested-reply [post]
func (ac *ArticleController) CreateNestedReply(c *gin.Context) {
	var reply domain.NestedReply

	replyID, err := strconv.Atoi(c.Param("reply-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&reply); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	reply.ReplyID = replyID
	reply.UserID = claims.ID
	reply.UserName = claims.Name

	replyResult, err := ac.arr.SaveNestedReply(reply)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, replyResult)
}

// DeleteNestedReply swagger
// @Summary delete nested reply
// @Accept json
// @Produce json
// @Param article-id path int true "article's id"
// @Param reply-id path int true "reply's id"
// @Param nested-reply-id path int true "nested reply's id"
// @Success 200 {object} integer "nested reply id"
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 403 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/articles/article/{article-id}/reply/{reply-id}/nested-reply/{nested-reply-id} [delete]
func (ac *ArticleController) DeleteNestedReply(c *gin.Context) {
	replyID, err := strconv.Atoi(c.Param("reply-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	nestedReplyID, err := strconv.Atoi(c.Param("nested-reply-id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	claims, exists := getAccessClaims(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "cannot find user claims"})
		return
	}

	nestedReply, err := ac.arr.FindNestedReplyByID(nestedReplyID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	if nestedReply.ReplyID != replyID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "nested reply's replyID is not mathed with param"})
		return
	}

	if nestedReply.UserID != claims.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "this nested reply is not yours"})
		return
	}

	if err := ac.arr.DeleteNestedReply(nestedReply); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nestedReplyID)
}

func getAccessClaims(c *gin.Context) (*domain.AccessClaim, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, exists
	}
	accessClaims := claims.(*domain.AccessClaim)
	return accessClaims, exists
}

// HttpError example for swagger
// not used
type HttpError struct {
	Detail string `json:"detail" example:"Some error comment"`
}
