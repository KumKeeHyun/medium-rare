package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/KumKeeHyun/medium-rare/article-service/config"
	"github.com/KumKeeHyun/medium-rare/article-service/controller"
	"github.com/KumKeeHyun/medium-rare/article-service/dao/sql"
	_ "github.com/KumKeeHyun/medium-rare/article-service/docs"
	"github.com/KumKeeHyun/medium-rare/article-service/middleware"
	"github.com/KumKeeHyun/medium-rare/article-service/util"

	"go.uber.org/zap"

	"github.com/chenjiandongx/ginprom"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Medium Rare Article Service
// @version 0.0.1

// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization

func main() {
	logger, err := util.BuildZapLogger()
	if err != nil {
		panic(err)
	}

	db, err := util.BuildMysqlConnection()
	if err != nil {
		panic(err)
	}

	syncProducer, err := util.BuildSyncProducer()
	// syncProducer, err := util.BuildMockSyncProducer()
	if err != nil {
		panic(err)
	}

	logger.Info("set dependency injection")

	arr := sql.NewSqlArticleReplyRepository(db)
	ac := controller.NewArticleController(arr, syncProducer, logger)

	logger.Info("set gin router")

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(ginprom.PromMiddleware(nil))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	jwtAuth := middleware.CheckJwtAuth()
	loggedIn := middleware.EnsureAuth()

	v1 := r.Group("/v1")
	{
		article := v1.Group("/articles")
		{
			article.GET("", ac.ListArticles)
			article.GET("/list", ac.ListArticlesByIDList)
			article.GET("/search", ac.SearchArticles)
			article.GET("/article/:article-id", jwtAuth, ac.GetArticle)
			article.POST("/article", jwtAuth, loggedIn, ac.CreateArticle)
			article.DELETE("/article/:article-id", jwtAuth, loggedIn, ac.DeleteArticle)

			reply := article.Group("/article/:article-id/reply", jwtAuth, loggedIn)
			{
				reply.POST("", ac.CreateReply)
				reply.DELETE("/:reply-id", ac.DeleteReply)

				nested := reply.Group("/:reply-id/nested-reply")
				{
					nested.POST("", ac.CreateNestedReply)
					nested.DELETE("/:nested-reply-id", ac.DeleteNestedReply)
				}
			}
		}
	}

	// for kubernetes liveness, readiness probe
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healty": "ok"})
	})

	logger.Info("start gin server",
		zap.String("addr", config.App.Address))

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}
