package service

import (
	"errors"
	"net/http"
	"personalBlog/model"
	"personalBlog/util"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PosetOp(r *gin.Engine) {
	type PostDetail struct {
		ID        uint
		CreatedAt time.Time
		UpdatedAt time.Time
		Title     string
		Content   string
		UserID    uint
	}
	// 创建文章
	r.POST("/createPost", util.MiddleWare(), func(c *gin.Context) {
		var post model.Post
		if err := c.ShouldBind(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		username, _ := c.Get("username")
		var user model.User
		if err := db.Omit("password").Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		post.UserID = user.ID
		post.User = user
		if err := db.Create(&post).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{
				"ID":       post.ID,
				"username": user.Username,
			},
			"msg": "success.",
		})
	})

	// 读取文章列表
	r.GET("/getPostList", util.MiddleWare(), func(c *gin.Context) {

		var storedPosts []model.Post
		if err := db.Find(&storedPosts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		var postLists []PostDetail
		for _, storedPost := range storedPosts {
			var post PostDetail
			post.ID = storedPost.ID
			post.CreatedAt = storedPost.CreatedAt
			post.UpdatedAt = storedPost.UpdatedAt
			post.Title = storedPost.Title
			post.Content = storedPost.Content
			post.UserID = storedPost.UserID
			postLists = append(postLists, post)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{
				"posts": postLists,
			},
			"total": len(storedPosts),
			"msg":   "success.",
		})
	})
	// 根据id获取文章
	r.GET("/getPost/:id", util.MiddleWare(), func(c *gin.Context) {
		id := c.Param("id")
		var storedPost model.Post
		if err := db.Find(&storedPost, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		var post PostDetail
		post.ID = storedPost.ID
		post.CreatedAt = storedPost.CreatedAt
		post.UpdatedAt = storedPost.UpdatedAt
		post.Title = storedPost.Title
		post.Content = storedPost.Content
		post.UserID = storedPost.UserID
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{
				"post": post,
			},
			"msg": "success.",
		})
	})
	// 更新文章
	r.POST("/updatePost", util.MiddleWare(), func(c *gin.Context) {
		postId := c.PostForm("id")
		// 判断该文章id存不存在
		var post model.Post
		err := db.Where("id = ?", postId).First(&post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": -1,
					"msg":  "Data does not exist.",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		// 判断是否文章作者
		username, _ := c.Get("username")
		var userOprate model.User
		if err := db.Omit("password").Where("username = ?", username).First(&userOprate).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		if post.UserID != userOprate.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  "Not the author of the article.",
			})
			return
		}

		var UpdatePost model.Post
		if err := c.ShouldBind(&UpdatePost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		// 更新
		if err := db.Debug().Model(&post).Where("id=?", postId).Updates(model.Post{Title: UpdatePost.Title, Content: UpdatePost.Content}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success.",
		})
	})

	// 删除文章
	r.DELETE("/deletePost/:id", util.MiddleWare(), func(c *gin.Context) {
		postId := c.Param("id")
		// 判断该文章id存不存在
		var post model.Post
		err := db.Where("id = ?", postId).First(&post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": -1,
					"msg":  "Data does not exist.",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}

		// 判断是否文章作者
		username, _ := c.Get("username")
		var userOperate model.User
		if err := db.Omit("password").Where("username = ?", username).First(&userOperate).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		if post.UserID != userOperate.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": -1,
				"msg":  "Not the author of the article.",
			})
			return
		}

		// 删除
		if err := db.Debug().Where("id=?", postId).Delete(&model.Post{}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  -1,
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success.",
		})
	})
}
