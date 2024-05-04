package dao

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"goPro/model"
)

var db2 *gorm.DB

func init() {
	var err error
	db2, err = gorm.Open("mysql", "root:123456@/mytest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db2.AutoMigrate(&model.Image{})
}

func AddImage(image model.Image) {
	db2.Create(&image)
}

func DeleteImage(image model.Image) int64 {

	res := db2.Where("username=?", image.Username).First(&image)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("没找到")
		return -1
		// 处理未找到记录的情况
	} else {
		//db.First(&post, id)

		res := db2.Delete(&image)
		return res.RowsAffected
	}
}

func SearchImage(username string) ([]model.Image, error) {
	var images []model.Image
	error := db2.Where("username=?", username).Find(&images).Error

	return images, error
}
