package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserName string `json:"username" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	//CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
	CreditCards []CreditCard `json:"creditcards" gorm:"many2many:user_creditcards;"`
}

type CreditCard struct {
	gorm.Model
	Number string `json:"number"  gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	UserID uint
}

func main() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Customer{}, &CreditCard{})

	// 第一个用户
	//db.Create(&Custom{
	//	UserName: "tom",
	//	CreditCards: []CreditCard{
	//		{Number: "12839054123151"},
	//		{Number: "12839054123152"},
	//		{Number: "12839054123153"},
	//	},
	//})
	// second user
	//db.Create(&Custom{
	//	UserName: "jason",
	//	CreditCards: []CreditCard{
	//		{Number: "12839054122345"},
	//		{Number: "12839054122346"},
	//		{Number: "12839054122347"},
	//	},
	//})

	//custome := &Custom{UserName: "tom"}
	//

	//var cs []CreditCard
	//
	//custom := &Custom{
	//	Model: gorm.Model{
	//		ID: 1,
	//	},
	//}

	//cards := []string{"12389451231253","12389451231456","123894512312343"}
	//db.Model(custom).Association("CreditCards").Find(&cs)
	//fmt.Println()

	var custom Customer
	db.Preload("CreditCards").Take(&custom, 1)
	fmt.Println(custom)
}
