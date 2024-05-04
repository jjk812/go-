package logic

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	"goPro/model"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@/mytest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{})
}

func InsertUser(user model.User) int {
	var result model.User
	h := md5.New()
	h.Write([]byte(user.Password))
	md5password := fmt.Sprintf("%x", h.Sum(nil))
	user.Password = md5password
	if err := db.Where("username = ?", user.Username).First(&result).Error; err != nil {
		// 用户名未被使用，可以创建新用户
		db.Create(&user)
		return 1

	} else {
		// 用户名已存在
		return -1
	}
}
func TestPassword(user model.User) bool {
	//从表单获取用户名和密码

	//使用MD5给密码加密
	h := md5.New()
	h.Write([]byte(user.Password))
	md5password := fmt.Sprintf("%x", h.Sum(nil))

	//查询数据库
	db.Where("username = ?", user.Username).First(&user)

	//检查密码
	if user.Password == md5password {
		return true
	} else {
		return false
	}
}
