package middleware

import (
	"fmt"
	"net/http"

	"github.com/KumKeeHyun/medium-rare/article-service/config"
	"github.com/KumKeeHyun/medium-rare/article-service/domain"
	jwt "github.com/dgrijalva/jwt-go"
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
