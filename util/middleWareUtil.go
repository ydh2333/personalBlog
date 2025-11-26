package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "The request header does not contain a token.",
			})
			c.Abort()
			return
		}
		parseToken, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		//fmt.Println("parseToken:", parseToken)
		c.Set("username", parseToken.Username)
		c.Next()

	}
}
