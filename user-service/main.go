package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/controller"
	"github.com/KumKeeHyun/medium-rare/user-service/dao/sql"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/KumKeeHyun/medium-rare/user-service/middleware"
	"github.com/KumKeeHyun/medium-rare/user-service/usecase"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := buildZapLogger()
	if err != nil {
		panic(err)
	}

	db, err := buildMysqlConnection()
	if err != nil {
		panic(err)
	}

	logger.Info("set dependency injection")

	// ur := memory.NewMemoryUserRepository()
	ur := sql.NewSqlUserRepository(db)
	uu := usecase.NewUserUsecase(ur, logger)
	au := usecase.NewAuthUsecase(ur, logger)
	uc := controller.NewUserController(uu, au, logger)

	logger.Info("set gin router")

	r := gin.Default()
	jwtAuth := middleware.JwtAuth()
	notLoggedIn := middleware.EnsureNotLoggedIn()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	v1 := r.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.GET("", uc.ListUsers)
			user.GET("/:id", uc.GetUser)
			user.POST("", notLoggedIn, uc.CreateUser)
			user.DELETE(":id", jwtAuth, uc.DeleteUser)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("", notLoggedIn, uc.Authorize)
			auth.POST("/refresh", uc.RefreshAuth)
		}
	}

	logger.Info("start gin server",
		zap.String("addr", config.App.Address))

	if err := r.Run(config.App.Address); err != nil {
		logger.Fatal("fail to start gin server",
			zap.Error(err))
	}
}

func buildZapLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(config.App.ZapConfig.Level)); err != nil {
		return nil, err
	}

	zapCfg := zap.Config{
		OutputPaths:       config.App.ZapConfig.OutputPaths,
		DisableCaller:     !config.App.ZapConfig.EableCaller,
		DisableStacktrace: !config.App.ZapConfig.EableCaller,
		Level:             level,
		Encoding:          config.App.ZapConfig.Encoding,
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
	}
	return zapCfg.Build()
}

func buildMysqlConnection() (*gorm.DB, error) {
	// temp url
	config.App.MysqlConfig.DbURL = "root:balns@tcp(192.168.219.204:3306)/userDB?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(config.App.MysqlConfig.DbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.User{})

	return db, nil
}
