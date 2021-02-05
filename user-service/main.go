package main

import (
	"io/ioutil"
	"time"

	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/controller"
	"github.com/KumKeeHyun/medium-rare/user-service/dao/sql"
	_ "github.com/KumKeeHyun/medium-rare/user-service/docs"
	"github.com/KumKeeHyun/medium-rare/user-service/middleware"
	"github.com/KumKeeHyun/medium-rare/user-service/usecase"
	"github.com/KumKeeHyun/medium-rare/user-service/util"
	"github.com/chenjiandongx/ginprom"
	ginzap "github.com/gin-contrib/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Medium Rare User Service
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

	// ur := memory.NewMemoryUserRepository()
	ur := sql.NewSqlUserRepository(db)
	uu := usecase.NewUserUsecase(ur, logger)
	au := usecase.NewAuthUsecase(ur, logger)
	uc := controller.NewUserController(uu, au, syncProducer, logger)

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
	notLoggedIn := middleware.EnsureNotAuth()

	v1 := r.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.GET("", uc.ListUsers)
			user.GET("/:id", jwtAuth, loggedIn, uc.GetUser)
			user.POST("", notLoggedIn, uc.CreateUser)
			user.DELETE(":id", jwtAuth, loggedIn, uc.DeleteUser)

			auth := user.Group("/auth")
			{
				auth.POST("", notLoggedIn, uc.Authorize)
				auth.POST("/refresh", uc.RefreshAuth)
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
