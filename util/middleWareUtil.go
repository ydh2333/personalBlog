package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrTokenFind)
			return
		}
		parseToken, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrTokenInvalid)
			return
		}
		//fmt.Println("parseToken:", parseToken)
		c.Set("username", parseToken.Username)
		c.Next()

	}
}
