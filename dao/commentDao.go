package dao

import "personalBlog/model"

func CreateComment(comm model.Comment) error {
	err := db.Create(&comm).Error
	return err
}

type CommOut struct {
	Content string
	UserID  uint
}

func FindCommentByPostId(postId string) ([]CommOut, error) {
	var comms []model.Comment
	err := db.Where("post_id = ?", postId).Find(&comms).Error
	var commOuts []CommOut
	for _, comm := range comms {
		var commOut CommOut
		commOut.Content = comm.Content
		commOut.UserID = comm.UserID
		commOuts = append(commOuts, commOut)
	}
	return commOuts, err
}
