package main

import (
	"log"

	"github.com/KumKeeHyun/medium-rare/user-service/controller"
	"github.com/KumKeeHyun/medium-rare/user-service/dao/memory"
	"github.com/KumKeeHyun/medium-rare/user-service/middleware"
	"github.com/KumKeeHyun/medium-rare/user-service/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	ur := memory.NewMemoryUserRepository()
	uu := usecase.NewUserUsecase(ur)
	au := usecase.NewAuthUsecase(ur)
	uc := controller.NewUserController(uu, au)

	r := gin.Default()
	jwtAuth := middleware.JwtAuth()
	notLoggedIn := middleware.EnsureNotLoggedIn()

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

	log.Fatal(r.Run(":8081"))
}
