package main

import (
	"personalBlog/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	service.RegisterUser(router)
	service.LoginUser(router)

	err := router.Run()
	if err != nil {
		panic(err)
	}

}
