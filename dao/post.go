package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"goPro/model"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@/mytest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Post{})
}

func AddPost(post model.Post) {
	db.Create(&post)
}

func DeletePost(post model.Post) int64 {

	res := db.Where("id = ?&& username=?", post.ID, post.Username).First(&post)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("没找到")
		return -1
		// 处理未找到记录的情况
	} else {
		//db.First(&post, id)
		fmt.Println("haha", post.ID)
		res := db.Delete(&post)
		return res.RowsAffected
	}
}

func UpdatePost(newPost model.Post) bool {
	var post model.Post
	fmt.Println("223", newPost.ID, newPost.Username)
	res := db.Where("id = ?&& username=?", newPost.ID, newPost.Username).First(&post)
	fmt.Println("2", post)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	} else {
		newPost.Title = post.Title
		newPost.ViewCount = post.ViewCount
		db.Model(&post).Update(newPost)
		return true
	}
}

func SearchPost(title string) ([]model.Post, error) {
	var posts []model.Post
	error := db.Where("title LIKE ?", "%"+title+"%").Find(&posts).Error

	return posts, error
}
func SearchUserPost(username string) ([]model.Post, error) {
	var posts []model.Post
	error := db.Where("username = ?", username).Find(&posts).Error
	return posts, error
}
func SearchPostID(id int) (model.Post, error) {
	var posts model.Post
	error := db.Where(" id=?", id).First(&posts).Error

	return posts, error
}
