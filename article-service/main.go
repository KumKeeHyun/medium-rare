package main

import (
	"github.com/KumKeeHyun/medium-rare/article-service/config"
	"github.com/KumKeeHyun/medium-rare/article-service/controller"
	"github.com/KumKeeHyun/medium-rare/article-service/dao/sql"
	"github.com/KumKeeHyun/medium-rare/article-service/middleware"
	"github.com/KumKeeHyun/medium-rare/article-service/util"
	"go.uber.org/zap"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger, err := util.BuildZapLogger()
	if err != nil {
		panic(err)
	}

	db, err := util.BuildMysqlConnection()
	if err != nil {
		panic(err)
	}

	// syncProducer, err := util.BuildSyncProducer()
	syncProducer, err := util.BuildMockSyncProducer()
	if err != nil {
		panic(err)
	}

	logger.Info("set dependency injection")

	arr := sql.NewSqlArticleReplyRepository(db)
	ac := controller.NewArticleController(arr, syncProducer, logger)

	logger.Info("set gin router")

	r := gin.Default()

	r.Use(ginprom.PromMiddleware(nil))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

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

	logger.Info("start gin server",
		zap.String("addr", config.App.Address))

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}
