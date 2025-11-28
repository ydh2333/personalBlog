package service

import (
	"errors"
	"fmt"
	"net/http"
	"personalBlog/dao"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommService(r *gin.Engine) {
	// 创建评论
	r.POST("/createComm", util.MiddleWare(), func(c *gin.Context) {
		var comm model.Comment
		if err := c.ShouldBind(&comm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		username, _ := c.Get("username")
		usernameStr, _ := username.(string)
		user, err := dao.GetUserByUsername(usernameStr)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		comm.UserID = user.ID

		if err := dao.CreateComment(comm); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})

	//读取评论
	r.GET("/getCommList/:postId", func(c *gin.Context) {
		// 判断文章是否存在
		postId := c.Param("postId")
		if _, err := dao.FindPostByID(postId); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				customErr := util.NewBusinessError(404, 2002, fmt.Sprintf("ID为%s的文章不存在", postId))
				_ = c.AbortWithError(http.StatusNotFound, customErr)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}
		// 查询该文章的所有评论
		commOuts, err := dao.FindCommentByPostId(postId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		c.JSON(http.StatusOK, util.Success(commOuts))

	})

}
