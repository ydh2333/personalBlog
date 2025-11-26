package service

import (
	"fmt"
	"net/http"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func CommOp(r *gin.Engine) {
	// 创建评论
	r.POST("/createComm", util.MiddleWare(), func(c *gin.Context) {
		var comm model.Comment
		if err := c.ShouldBind(&comm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		value, _ := c.Get("username")
		var userID uint
		if err := db.Debug().Model(&model.User{}).Select("id").Where("username = ?", value).Scan(&userID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		comm.UserID = userID
		fmt.Println("comm:", comm)
		if err := db.Create(&comm).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success.",
		})
	})
	//读取评论
	r.GET("/getCommList/:postId", func(c *gin.Context) {
		type CommOut struct {
			Content string
			UserID  uint
		}
		postId := c.Param("postId")
		var comms []model.Comment
		if err := db.Debug().Where("post_id = ?", postId).Find(&comms).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var commOuts []CommOut
		for _, comm := range comms {
			var commOut CommOut
			commOut.Content = comm.Content
			commOut.UserID = comm.UserID
			commOuts = append(commOuts, commOut)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": commOuts,
			"msg":  "success.",
		})

	})

}
