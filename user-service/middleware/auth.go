package middleware

import (
	"fmt"
	"net/http"

	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CheckJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tStr := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tStr, &domain.AccessClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.App.JWTSecret), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*domain.AccessClaim); ok && token.Valid {
				c.Set("claims", claims)
			}
		}

		c.Next()
	}
}

func EnsureAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exist := c.Get("claims"); !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "authentication is required"})
			return
		}
		c.Next()
	}
}

func EnsureNotAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		if _, exists := c.Get("claims"); exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "authentication is not required"})
			return
		}

		c.Next()
	}
}
