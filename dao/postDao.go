package dao

import (
	"personalBlog/model"
	"time"
)

func CreatePost(post model.Post) error {
	err := db.Create(&post).Error
	return err
}

type PostDetail struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Content   string
	UserID    uint
}

func FindPostAll() ([]PostDetail, error) {

	var storedPosts []model.Post
	err := db.Find(&storedPosts).Error

	var postLists []PostDetail
	for _, storedPost := range storedPosts {
		var post PostDetail
		post.ID = storedPost.ID
		post.CreatedAt = storedPost.CreatedAt
		post.UpdatedAt = storedPost.UpdatedAt
		post.Title = storedPost.Title
		post.Content = storedPost.Content
		post.UserID = storedPost.UserID
		postLists = append(postLists, post)
	}

	return postLists, err
}

func FindPostByID(id string) (PostDetail, error) {
	var storedPost model.Post
	err := db.Where("id=?", id).First(&storedPost).Error

	var post PostDetail
	post.ID = storedPost.ID
	post.CreatedAt = storedPost.CreatedAt
	post.UpdatedAt = storedPost.UpdatedAt
	post.Title = storedPost.Title
	post.Content = storedPost.Content
	post.UserID = storedPost.UserID

	return post, err
}

func UpdatePost(postId string, UpdatePost model.Post) error {

	err := db.Debug().Model(&UpdatePost).Where("id=?", postId).Updates(model.Post{Title: UpdatePost.Title, Content: UpdatePost.Content}).Error
	return err
}

func DeletePost(id string) error {
	err := db.Debug().Where("id=?", id).Delete(&model.Post{}).Error

	return err
}
