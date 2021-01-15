package main

import (
	"github.com/KumKeeHyun/medium-rare/reading-list-service/config"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/controller"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/dao/sql"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/middleware"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/util"
	"github.com/KumKeeHyun/medium-rare/reading-list-service/util/erouter"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	rr := sql.NewSqlReadingRepository(db)
	ec := controller.NewEventController(rr, logger)
	rc := controller.NewReadingController(rr, logger)

	er := erouter.NewEventRouter("reading-list", logger)
	er.SetHandler("read-article", ec.ReadArticle)
	if err := er.StartRouter(); err != nil {
		panic(err)
	}

	r := gin.Default()

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

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}
