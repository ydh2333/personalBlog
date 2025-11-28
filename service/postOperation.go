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

func PosetOp(r *gin.Engine) {

	// 创建文章
	r.POST("/createPost", util.MiddleWare(), func(c *gin.Context) {
		var post model.Post
		if err := c.ShouldBind(&post); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		username, _ := c.Get("username")
		usernameStr, _ := username.(string)
		storedUser, err := dao.GetUserByUsername(usernameStr)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		post.UserID = storedUser.ID

		if err := dao.CreatePost(post); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		c.JSON(http.StatusOK, util.Success(nil))
	})

	// 读取文章列表
	r.GET("/getPostList", util.MiddleWare(), func(c *gin.Context) {

		postLists, err := dao.FindPostAll()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}

		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"data": gin.H{
					"posts": postLists,
				},
				"total": len(postLists),
			}))
	})
	// 根据id获取文章
	r.GET("/getPost/:id", util.MiddleWare(), func(c *gin.Context) {
		id := c.Param("id")

		post, err := dao.FindPostByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				customErr := util.NewBusinessError(404, 2002, fmt.Sprintf("ID为%s的文章不存在", id))
				_ = c.AbortWithError(http.StatusNotFound, customErr)
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			}
			return
		}

		c.JSON(http.StatusOK, util.Success(post))
	})
	// 更新文章
	r.POST("/updatePost", util.MiddleWare(), func(c *gin.Context) {
		postId := c.PostForm("id")
		// 判断该文章id存不存在
		post, err := dao.FindPostByID(postId)
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
		usernameStr, _ := username.(string)
		storedUser, err := dao.GetUserByUsername(usernameStr)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		if post.UserID != storedUser.ID {
			c.JSON(http.StatusBadRequest, util.NewBusinessError(400, 3002, "非作者本人，无法更新文章"))
			return
		}

		var UpdatePost model.Post
		if err := c.ShouldBind(&UpdatePost); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		// 更新
		if err := dao.UpdatePost(postId, UpdatePost); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})

	// 删除文章
	r.DELETE("/deletePost/:id", util.MiddleWare(), func(c *gin.Context) {
		postId := c.PostForm("id")
		// 判断该文章id存不存在
		post, err := dao.FindPostByID(postId)
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
		usernameStr, _ := username.(string)
		storedUser, err := dao.GetUserByUsername(usernameStr)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		if post.UserID != storedUser.ID {
			c.JSON(http.StatusBadRequest, util.NewBusinessError(400, 3002, "非作者本人，无法更新文章"))
			return
		}

		// 删除
		if err := dao.DeletePost(postId); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		c.JSON(http.StatusOK, util.Success(nil))
	})
}
