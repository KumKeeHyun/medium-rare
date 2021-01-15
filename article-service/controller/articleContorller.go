package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// GET /v1/articles?p=
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

// GET /v1/articles/list
func (ac *ArticleController) ListArticlesByIDList(c *gin.Context) {
	type idListRequestBody struct {
		IDs []int `json:"ids"`
	}
	var reqBody idListRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	articles, err := ac.arr.FindArticleByIDList(reqBody.IDs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article_list": articles})
}

// GET /v1/articles/search?q=
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

// GET /v1/articles/:article-id
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
		}

		mshld, _ := json.Marshal(rae)
		ac.producer.SendMessage(&sarama.ProducerMessage{
			Topic: "read-article",
			Value: sarama.ByteEncoder(mshld),
		})
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// POST /v1/articles
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

	c.JSON(http.StatusOK, gin.H{"created": articleResult})
}

// DELETE /v1/articles/:article-id
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

// POST /v1/articles/:article-id/reply
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

	c.JSON(http.StatusOK, gin.H{"reply": replyResult})
}

// DELETE /v1/articles/:article-id/reply/:reply-id
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

// POST /v1/articles/:article-id/reply/:reply-id/nested-reply
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

	c.JSON(http.StatusOK, gin.H{"reply": replyResult})
}

// POST /v1/articles/:article-id/reply/:reply-id/nested-reply/:nested-reply-id
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