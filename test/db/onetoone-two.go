package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type Blog struct {
	//gorm.Model
	ID         int    `json:"id" gorm:"autoIncrement;primaryKey"`
	Title      string `gorm:"not null;size:256"`
	Content    string `gorm:"type:text;not null"`
	BlogType   BlogType
	BlogTypeID int `gorm:"type:int;not null"`
}

type BlogType struct {
	ID   int `json:"id" gorm:"autoIncrement;primaryKey"`
	Name string
}

func main() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: newLogger,
	})
	db.AutoMigrate(&Blog{}, &BlogType{})

	//db.Create(&Blog{
	//	Title:   "golang study",
	//	Content: "i love golang",
	//	BlogType: BlogType{
	//		Name: "original",
	//	},
	//})
	//db.Create(&Blog{
	//	Title:   "python study",
	//	Content: "i love python",
	//	BlogType: BlogType{
	//		Name: "reproduce",
	//	},
	//})
	//
	db.Create(&Blog{
		Title:   "java study",
		Content: "i love java",
		BlogType: BlogType{
			Name: "reproduce",
		},
	})

	db.Create(&Blog{
		Title:   "c study",
		Content: "i love c",
		BlogType: BlogType{
			Name: "reproduce",
		},
	})

	var blogs []Blog
	db.Find(&blogs)
	fmt.Println(blogs)

	var blog Blog
	db.Where("title = ?", "golang study").Find(&blog)
	fmt.Println(blog)

	db.Where("id = ?", "2").Find(&blog)
	fmt.Println(blog)

	db.First(&blog, 2)
	fmt.Println(blog)

	db.Last(&blogs)
	fmt.Println(blogs)
}
