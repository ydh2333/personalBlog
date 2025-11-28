package dao

import (
	"personalBlog/model"
	"personalBlog/util"
)

var db, _ = util.SqlConnect()

func GetUserByUsername(username string) (model.User, error) {
	var storedUser model.User
	err := db.Where("username = ?", username).First(&storedUser).Error
	return storedUser, err
}

func InsertUser(user model.User) error {
	err := db.Create(&user).Error
	return err
}

// 对外暴露脱敏信息
type OutUser struct {
	ID       uint   `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func SelectUserAll() ([]OutUser, error) {

	var users []model.User
	var outUsers []OutUser
	err := db.Omit("password").Find(&users).Error
	for _, user := range users {
		var outUser OutUser
		outUser.ID = user.ID
		outUser.Username = user.Username
		outUser.Email = user.Email
		outUsers = append(outUsers, outUser)
	}
	return outUsers, err
}
