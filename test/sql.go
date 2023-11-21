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

//type User struct {
//	Id       int
//	Username string
//	Orders   []Order
//}

type Order struct {
	Id     int
	UserID uint
	Price  float64
}

func main() {
	// 数据库连接
	dsn := "root:123456@(192.168.10.168:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// 初始化表
	//_ = db.AutoMigrate(User{}, Order{})

	// 创建数据
	//_ = CreatUser(db)
	//_ = CreatOrder(db)
	UserOrder(db)
	//fmt.Println(GetUser(db))
}

func CreatUser(db *gorm.DB) error {
	err := db.Create(&[]User{
		{Username: "little_A"},
		{Username: "little_B"},
		{Username: "little_C"},
		{Username: "little_D"},
		{Username: "little_H"},
		{Username: "little_E"},
	}).Error
	return err
}

func CreatOrder(db *gorm.DB) error {
	err := db.Create(&[]Order{
		{UserID: 1, Price: 1},
		{UserID: 1, Price: 2},
		{UserID: 1, Price: 3},
		{UserID: 1, Price: 4},
		{UserID: 2, Price: 5},
		{UserID: 2, Price: 6},
		{UserID: 2, Price: 7},
		{UserID: 3, Price: 8},
		{UserID: 3, Price: 9},
		{UserID: 4, Price: 10},
		{UserID: 5, Price: 10},
		{UserID: 6, Price: 10},
	}).Error
	return err
}

func GetUser(db *gorm.DB) ([]User, error) {
	users := make([]User, 0)
	if err := db.Preload("Orders").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UserOrder(db *gorm.DB) error {
	users := make([]User, 0)
	if err := db.Order("username").Find(&users).Error; err != nil {
		return err
	}
	fmt.Println(users)
	return nil
}
