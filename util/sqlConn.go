package util

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SqlConnect() (*gorm.DB, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/personal_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, ErrDBConnect
	}
	return db, nil
}
