package service

import (
	"net/http"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func UserOp(r *gin.Engine) {
	// 对外暴露脱敏信息
	type OutUser struct {
		ID       uint   `json:"ID"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	// 获取所有用的信息
	r.GET("/userAll", util.MiddleWare(), func(c *gin.Context) {
		var users []model.User
		var outUsers []OutUser
		db.Omit("password").Find(&users)
		for _, user := range users {
			var outuser OutUser
			outuser.ID = user.ID
			outuser.Username = user.Username
			outuser.Email = user.Email
			outUsers = append(outUsers, outuser)
		}
		// 获取操作人
		value, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"operator": value,
			"data":     outUsers,
		})
	})
}
