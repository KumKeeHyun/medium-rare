package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KumKeeHyun/medium-rare/user-service/adapter"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/KumKeeHyun/medium-rare/user-service/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	uu  usecase.UserUsecase
	au  usecase.AuthUsecase
	log *zap.Logger
}

func NewUserController(uu usecase.UserUsecase, au usecase.AuthUsecase, log *zap.Logger) *UserController {
	return &UserController{
		uu:  uu,
		au:  au,
		log: log,
	}
}

// ListUsers ...
// GET
// /users
func (uc *UserController) ListUsers(c *gin.Context) {
	us, err := uc.uu.FindUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	aus := adapter.ToAdapterUsers(us)
	c.JSON(http.StatusOK, aus)
}

// GetUser ...
// GET
// /users/:id
func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := uc.uu.FindOne(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adapter.ToAdapterUser(&u))
}

// CreateUser ...
// POST
// /users
// Body : json(domain.User:Email,Password,Name,Gender,Birth)
func (uc *UserController) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, err := uc.uu.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adapter.ToAdapterUser(&userResult))
}

// DeleteUser ...
// DELETE
// /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	claims, err := getAccessClaime(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if claims.ID != id {
		c.AbortWithStatus(http.StatusForbidden)
	}

	if err := uc.uu.Unregister(domain.User{ID: id}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, id)
}

// Authorize ...
// POST
// /auth
// Body : json(domain.User:Email,Password)
func (uc *UserController) Authorize(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := uc.au.Login(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, tokenPair)
}

// RefreshAuth ...
// POST
// /auth/refresh
// Body : json(refreshToken)
func (uc *UserController) RefreshAuth(c *gin.Context) {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	var tokenReq tokenReqBody
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := uc.au.RefreshToken(tokenReq.RefreshToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
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
