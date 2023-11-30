package main

import (
	"app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/testuser?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&model.Like{}, &model.User{}, &model.Comment{}, &model.Article{}, &model.Category{}, &model.Tag{}, &model.UserInfo{})
}
