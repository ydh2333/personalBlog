package main

import (
	"personalBlog/model"
	"personalBlog/util"
)

func main() {
	var db, _ = util.SqlConnect()
	err := db.Debug().AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		panic(err)
	}
}
