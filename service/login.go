package service

import (
	"net/http"
	"personalBlog/model"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var db = util.SqlConnect()

func RegisterUser(r *gin.Engine) {
	r.POST("/register", func(c *gin.Context) {
		user := model.User{}
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
		// 入库
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{
				"ID":       user.ID,
				"Username": user.Username,
				"email":    user.Email,
			},
			"message": "success.",
		})
	})
}

func LoginUser(r *gin.Engine) {
	r.POST("/login", func(c *gin.Context) {
		user := model.User{}
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 查询用户名是否存在
		var storedUser model.User
		if err = db.Debug().Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The username or password is incorrect.\n"})
			return
		}
		// 验证密码是否正确
		if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The username or password is incorrect.\n"})
			return
		}
		// 生成jwt
		token := util.GenerateToken(user.Username)

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{
				"ID":       storedUser.ID,
				"Username": storedUser.Username,
				"token":    token,
			},
			"message": "success.",
		})
	})

}
