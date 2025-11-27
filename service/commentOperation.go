package service

import (
	"errors"
	"fmt"
	"net/http"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommOp(r *gin.Engine) {
	// 创建评论
	r.POST("/createComm", util.MiddleWare(), func(c *gin.Context) {
		var comm model.Comment
		if err := c.ShouldBind(&comm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		value, _ := c.Get("username")
		var userID uint
		if err := db.Model(&model.User{}).Select("id").Where("username = ?", value).Scan(&userID).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		comm.UserID = userID
		fmt.Println("comm:", comm)
		if err := db.Create(&comm).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})
	//读取评论
	r.GET("/getCommList/:postId", func(c *gin.Context) {
		type CommOut struct {
			Content string
			UserID  uint
		}
		// 判断文章是否存在
		postId := c.Param("postId")
		if err := db.Debug().Where("id = ?", postId).First(&model.Post{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				customErr := util.NewBusinessError(404, 2002, fmt.Sprintf("ID为%s的文章不存在", postId))
				_ = c.AbortWithError(http.StatusNotFound, customErr)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}
		// 查询该文章的所有评论
		var comms []model.Comment
		if err := db.Debug().Where("post_id = ?", postId).Find(&comms).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		var commOuts []CommOut
		for _, comm := range comms {
			var commOut CommOut
			commOut.Content = comm.Content
			commOut.UserID = comm.UserID
			commOuts = append(commOuts, commOut)
		}
		c.JSON(http.StatusOK, util.Success(commOuts))

	})

}
