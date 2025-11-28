package main

import (
	"net/http"
	"personalBlog/service"
	"personalBlog/util"

	"github.com/gin-gonic/gin"
)

func main() {

	config := util.Config{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      true,
		FileLoggingEnabled:    true,
		Directory:             "./logs",
		Filename:              "./logs/",
		MaxSize:               1,
		MaxBackups:            10,
		MaxAge:                30,
		Level:                 1,
	}

	util.InitLogger(config)

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
	service.UserOp(router)
	service.PosetOp(router)
	service.CommOp(router)

	err := router.Run()
	if err != nil {
		panic(err)
	}

}
