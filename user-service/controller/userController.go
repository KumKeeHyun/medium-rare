package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KumKeeHyun/medium-rare/user-service/adapter"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/KumKeeHyun/medium-rare/user-service/usecase"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	uu       usecase.UserUsecase
	au       usecase.AuthUsecase
	producer sarama.SyncProducer
	log      *zap.Logger
}

func NewUserController(uu usecase.UserUsecase, au usecase.AuthUsecase, sp sarama.SyncProducer, log *zap.Logger) *UserController {
	return &UserController{
		uu:       uu,
		au:       au,
		producer: sp,
		log:      log,
	}
}

// ListUsers swagger
// @Summary List all users
// @Accept json
// @Produce json
// @Success 200 {object} adapter.User
// @Failure 500 {object} controller.HttpError
// @Router /v1/users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	us, err := uc.uu.FindUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	aus := adapter.ToAdapterUsers(us)
	c.JSON(http.StatusOK, aus)
}

// GetUser swagger
// @Summary Show user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} adapter.User
// @Failure 500 {object} controller.HttpError
// @Router /v1/users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	claims, err := getAccessClaime(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "unauthorized user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if claims.ID != id {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "URL param is not user's id"})
		return
	}

	u, err := uc.uu.FindOne(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adapter.ToAdapterUser(&u))
}

// CreateUser swagger
// @Summary create user
// @Accept json
// @Produce json
// @Param user body domain.CreateUser true "user info"
// @Success 200 {object} adapter.User
// @Failure 400 {object} controller.HttpError "fail to unmashal body"
// @Failure 500 {object} controller.HttpError
// @Router /v1/users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	userResult, err := uc.uu.Register(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	mshld, _ := json.Marshal(adapter.ToAdapterCreateUserEvent(&userResult))
	uc.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "create-user",
		Value: sarama.ByteEncoder(mshld),
	})

	c.JSON(http.StatusOK, adapter.ToAdapterUser(&userResult))
}

// DeleteUser swagger
// @Summary delete user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 int int "user id"
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Failure 403 {object} controller.HttpError
// @Failure 500 {object} controller.HttpError
// @Security JWTToken
// @Router /v1/users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	claims, err := getAccessClaime(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "unauthorized user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if claims.ID != id {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"detail": "URL param is not user's id"})
		return
	}

	user := domain.User{ID: id}
	if err := uc.uu.Unregister(user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	mshld, _ := json.Marshal(adapter.ToAdapterDeleteUserEvent(&user))
	uc.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "delete-user",
		Value: sarama.ByteEncoder(mshld),
	})

	c.JSON(http.StatusOK, id)
}

// Authorize swagger
// @Summary authorize user
// @Accept json
// @Produce json
// @Param user body domain.LoginUser true "user info"
// @Success 200 {object} domain.TokenPair
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Router /v1/users/auth [post]
func (uc *UserController) Authorize(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	tokenPair, err := uc.au.Login(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// RefreshAuth swagger
// @Summary refresh jwt token
// @Accept json
// @Produce json
// @Param refresh_token body domain.RefreshToken true "refresh token"
// @Success 200 {object} domain.AccessToken
// @Failure 400 {object} controller.HttpError
// @Failure 401 {object} controller.HttpError
// @Router /v1/users/auth/refresh [post]
func (uc *UserController) RefreshAuth(c *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	var tokenReq tokenReqBody
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	accessToken, err := uc.au.RefreshToken(tokenReq.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func getAccessClaime(c *gin.Context) (*domain.AccessClaim, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, fmt.Errorf("claims must set in context")
	}

	accessClaims := claims.(*domain.AccessClaim)
	return accessClaims, nil
}

// HttpError example for swagger
// not used
type HttpError struct {
	Detail string `json:"detail" example:"Some error comment"`
}
