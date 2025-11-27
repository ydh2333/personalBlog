package service

import (
	"errors"
	"fmt"
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
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		username, _ := c.Get("username")
		var user model.User
		if err := db.Omit("password").Where("username = ?", username).First(&user).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		post.UserID = user.ID
		if err := db.Create(&post).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"ID":       post.ID,
				"username": user.Username,
			}))
	})

	// 读取文章列表
	r.GET("/getPostList", util.MiddleWare(), func(c *gin.Context) {

		var storedPosts []model.Post
		if err := db.Find(&storedPosts).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
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

		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"data": gin.H{
					"posts": postLists,
				},
				"total": len(storedPosts),
			}))
	})
	// 根据id获取文章
	r.GET("/getPost/:id", util.MiddleWare(), func(c *gin.Context) {
		id := c.Param("id")
		var storedPost model.Post
		if err := db.Where("id=?", id).First(&storedPost).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				customErr := util.NewBusinessError(404, 2002, fmt.Sprintf("ID为%s的文章不存在", id))
				_ = c.AbortWithError(http.StatusNotFound, customErr)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}
		var post PostDetail
		post.ID = storedPost.ID
		post.CreatedAt = storedPost.CreatedAt
		post.UpdatedAt = storedPost.UpdatedAt
		post.Title = storedPost.Title
		post.Content = storedPost.Content
		post.UserID = storedPost.UserID
		c.JSON(http.StatusOK, util.Success(post))
	})
	// 更新文章
	r.POST("/updatePost", util.MiddleWare(), func(c *gin.Context) {
		postId := c.PostForm("id")
		// 判断该文章id存不存在
		var post model.Post
		err := db.Where("id = ?", postId).First(&post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.AbortWithError(http.StatusNotFound, util.ErrArticleNotFound)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}

		// 判断是否文章作者
		username, _ := c.Get("username")
		var userOprate model.User
		if err := db.Omit("password").Where("username = ?", username).First(&userOprate).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		if post.UserID != userOprate.ID {
			c.JSON(http.StatusBadRequest, util.NewBusinessError(400, 3002, "非作者本人，无法更新文章"))
			return
		}

		var UpdatePost model.Post
		if err := c.ShouldBind(&UpdatePost); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		// 更新
		if err := db.Debug().Model(&post).Where("id=?", postId).Updates(model.Post{Title: UpdatePost.Title, Content: UpdatePost.Content}).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})

	// 删除文章
	r.DELETE("/deletePost/:id", util.MiddleWare(), func(c *gin.Context) {
		postId := c.Param("id")
		// 判断该文章id存不存在
		var post model.Post
		err := db.Where("id = ?", postId).First(&post).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.AbortWithError(http.StatusNotFound, util.ErrArticleNotFound)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}

		// 判断是否文章作者
		username, _ := c.Get("username")
		var userOperate model.User
		if err := db.Omit("password").Where("username = ?", username).First(&userOperate).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		if post.UserID != userOperate.ID {
			c.JSON(http.StatusBadRequest, util.NewBusinessError(400, 3002, "非作者本人，无法更新文章"))
			return
		}

		// 删除
		if err := db.Debug().Where("id=?", postId).Delete(&model.Post{}).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})
}
