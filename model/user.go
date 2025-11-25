package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" form:"username" json:"username" binding:"required"`
	Password string `gorm:"not null" form:"password" json:"password" binding:"required"`
	Email    string `gorm:"unique;not null" form:"email" json:"email"`
}
