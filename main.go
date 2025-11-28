package main

import (
	"net/http"
	"personalBlog/service"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func main() {

	util.InitLogger()

	router := gin.New()
	router.Use(util.GinLogger(), util.GinRecovery(true))

	//router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, util.ErrResourceNotFound)
	})
	// 统一错误处理中间件
	router.Use(util.ErrorHandler())
	service.RegisterUser(router)
	service.LoginUser(router)
	service.UserService(router)
	service.PosetService(router)
	service.CommService(router)

	err := router.Run()
	if err != nil {
		panic(err)
	}
}
