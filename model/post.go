package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" form:"title" json:"title" binding:"required"`
	Content string `gorm:"not null" form:"content" json:"content" binding:"required"`
	UserID  uint   `binding:"-"`
	User    User   `binding:"-"`
}
