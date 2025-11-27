package service

import (
	"net/http"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var db, _ = util.SqlConnect()

func RegisterUser(r *gin.Engine) {
	// 注册
	r.POST("/register", func(c *gin.Context) {
		user := model.User{}
		err := c.ShouldBind(&user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.NewBusinessError(http.StatusInternalServerError, 5003, "哈希加密失败"))
			return
		}
		user.Password = string(hashedPassword)
		// 入库
		if err := db.Create(&user).Error; err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, util.ErrSystemError)
			return
		}
		// 返回结果
		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"ID":       user.ID,
				"Username": user.Username,
				"email":    user.Email,
			}))
	})
}

func LoginUser(r *gin.Engine) {
	r.POST("/login", func(c *gin.Context) {
		user := model.User{}
		err := c.ShouldBind(&user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, util.ErrInvalidParam)
			return
		}

		// 查询用户名是否存在
		var storedUser model.User
		if err = db.Debug().Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
			_ = c.AbortWithError(http.StatusNotFound, util.NewBusinessError(404, 2004, "用户名或密码不正确"))
			return
		}
		// 验证密码是否正确
		if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
			_ = c.AbortWithError(http.StatusNotFound, util.NewBusinessError(404, 2004, "用户名或密码不正确"))
			return
		}
		// 生成jwt
		token, err := util.GenerateToken(user.Username)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, util.ErrAuthFailed)
			return
		}

		c.JSON(http.StatusOK, util.Success(
			gin.H{
				"ID":       storedUser.ID,
				"Username": storedUser.Username,
				"token":    token,
			}))
	})

}
