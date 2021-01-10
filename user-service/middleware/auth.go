package middleware

import (
	"fmt"
	"net/http"

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

			return []byte("kkh"), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusProxyAuthRequired, gin.H{"error": err.Error()})
		}

		claims, ok := token.Claims.(*domain.AccessClaim)
		if !ok {
			c.AbortWithStatusJSON(http.StatusNonAuthoritativeInfo, gin.H{"error-token": token.Claims})
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("claims")
		if exists {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
