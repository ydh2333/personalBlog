package service

import (
	"net/http"
	"personalBlog/dao"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func UserService(r *gin.Engine) {
	// 获取所有用户的信息
	r.GET("/userAll", util.MiddleWare(), func(c *gin.Context) {

		userAll, err := dao.SelectUserAll()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		// 获取操作人
		value, _ := c.Get("username")
		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"operator": value,
				"data":     userAll,
			}))
	})
}
