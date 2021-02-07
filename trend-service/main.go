package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/KumKeeHyun/medium-rare/trend-service/config"
	"github.com/KumKeeHyun/medium-rare/trend-service/controller"
	"github.com/KumKeeHyun/medium-rare/trend-service/dao/sql"
	_ "github.com/KumKeeHyun/medium-rare/trend-service/docs"
	"github.com/KumKeeHyun/medium-rare/trend-service/middleware"
	"github.com/KumKeeHyun/medium-rare/trend-service/util"
	"github.com/KumKeeHyun/medium-rare/trend-service/util/erouter"
	"github.com/go-resty/resty/v2"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Medium Rare Trend Service
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

	logger.Info("set dependency injection")

	rrr := sql.NewSqlReadRecordRepository(db)
	ec := controller.NewEventController(rrr, logger)
	tc := controller.NewTrendController(rrr, logger)

	er := erouter.NewEventRouter("trend", logger)
	er.SetHandler("read-article", ec.ReadArticle)

	if err := er.StartRouter(); err != nil {
		panic(err)
	}
	defer er.Stop()

	logger.Info("set gin router")

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	jwtAuth := middleware.CheckJwtAuth()
	loggedIn := middleware.EnsureAuth()

	v1 := r.Group("/v1")
	{
		trend := v1.Group("/trend")
		{
			trend.GET("", tc.ListTrend)
			trend.GET("/user", jwtAuth, loggedIn, tc.ListTrendForUser)
		}
	}

	logger.Info("start gin server",
		zap.String("addr", config.App.Address))

	// for kubernetes liveness probe
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healty": "ok"})
	})

	// for kubernetes readiness probe
	r.GET("/ready", func(c *gin.Context) {
		url := fmt.Sprintf("http://%s:%s/healthy", config.App.ArticleConfig.Host, config.App.ArticleConfig.Port)
		resp, _ := resty.New().R().SetHeader("Content-Type", "application/json").Get(url)

		if resp.IsSuccess() {
			c.JSON(http.StatusOK, gin.H{"ready": "ok"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"ready": "not ok"})
		}
	})

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}
