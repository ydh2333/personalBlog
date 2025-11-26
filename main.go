package main

import (
	"personalBlog/service"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	// 统一错误处理中间件
	router.Use(util.ErrorHandler())
	service.RegisterUser(router)
	service.LoginUser(router)
	service.UserOp(router)
	service.PosetOp(router)
	service.CommOp(router)

	err := router.Run()
	if err != nil {
		panic(err)
	}

}
