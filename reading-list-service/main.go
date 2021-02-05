package main

import (
	"io/ioutil"
	"time"

	"github.com/KumKeeHyun/medium-rare/reading-list-service/config"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/controller"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/dao/sql"
	_ "github.com/KumKeeHyun/medium-rare/reading-list-service/docs"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/middleware"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/util"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/util/erouter"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

// @title Medium Rare Reading List Service
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

	rr := sql.NewSqlReadingRepository(db)
	ec := controller.NewEventController(rr, logger)
	rc := controller.NewReadingController(rr, logger)

	er := erouter.NewEventRouter("reading-list", logger)
	er.SetHandler("read-article", ec.ReadArticle)
	if err := er.StartRouter(); err != nil {
		panic(err)
	}

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
		list := v1.Group("/reading-list", jwtAuth, loggedIn)
		{
			list.GET("/recent", rc.ListRecent)
			list.POST("/saved", rc.SaveArticle)
			list.GET("/saved", rc.ListSaved)
		}
	}

	logger.Info("start gin server",
		zap.String("addr", config.App.Address))

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}
