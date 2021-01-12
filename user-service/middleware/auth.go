package middleware

import (
	"fmt"
	"net/http"

	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tStr := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tStr, &domain.AccessClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.App.JWTSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": err.Error()})
			return
		}

		claims, ok := token.Claims.(*domain.AccessClaim)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "Unexpected claim type or Unvaild token"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("claims")
		if exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "this endpoint don't require authorization"})
			return
		}

		c.Next()
	}
}
