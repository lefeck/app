package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	UserName string `json:"username" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	Class    Class  `json:"class" gorm:"foreignKey:ClassID"`
	ClassID  uint
}

type Class struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
}

func main() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&Student{}, &Class{})

	//db.Create(&Student{
	//	UserName: "john",
	//	Class: Class{
	//		Name: "eight class",
	//	},
	//})
	//db.Create(&Student{
	//	UserName: "lack",
	//	Class: Class{
	//		Name: "nine class",
	//	},
	//})

	var stu Student
	db.Find(&stu)
	fmt.Println(stu)
}
