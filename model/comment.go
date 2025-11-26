package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"commContent" form:"commContent" binding:"required"`
	UserID  uint   `gorm:"not null" binding:"-"`
	User    User   `binding:"-"`
	PostID  uint   `gorm:"not null" json:"postId" form:"postId" binding:"required"`
	Post    Post   `binding:"-"`
}
