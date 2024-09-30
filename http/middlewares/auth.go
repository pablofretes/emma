package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	utils "emma/internal/utils"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		tokenSplitted := strings.Split(authorization, " ")
		if len(tokenSplitted) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "",
				"message": "unauthorized",
				"success": false,
			})
			return
		}

		token, err := utils.Authenticate(tokenSplitted[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   err.Error(),
				"message": "unauthorized",
				"success": false,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   claims.Valid().Error(),
				"message": "unauthorized",
				"success": false,
			})
			return
		}

		role := claims["role"].(string)
		roleSlice := []string{role}

		if utils.SomeElementInSlice(roleSlice, roles) {
			c.Set("user", claims)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":   "",
			"message": "forbidden",
			"success": false,
		})

	}
}
