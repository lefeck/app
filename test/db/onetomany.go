package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*

// 第一种方式:
type Customer struct {
	gorm.Model
	UserName    string       `json:"username" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string `json:"number"  gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	CustomerID uint
}
*/

// 第二种方式:
type Customer struct {
	gorm.Model
	UserName    string       `json:"username" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	CreditCards []CreditCard `gorm:"foreignKey:UserID"`
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

	//// fist create
	//db.Create(&Customer{
	//	UserName: "tom",
	//	CreditCards: []CreditCard{
	//		{Number: "12839054123151"},
	//		{Number: "12839054123152"},
	//		{Number: "12839054123153"},
	//	},
	//})
	////second create
	//db.Create(&Customer{
	//	UserName: "jason",
	//	CreditCards: []CreditCard{
	//		{Number: "12839054122345"},
	//		{Number: "12839054122346"},
	//		{Number: "12839054122347"},
	//	},
	//})

	var custom Customer
	// 查询功能（预加载）
	db.Model(&custom).Preload("CreditCards").First(&custom, 1)
	//下面的方法等同于上面
	//db.Model(&custom).Preload("CreditCards").Take(&custom, 1)
	fmt.Println(custom)
	//
	//// preload
	//var user []Customer //设置接收的切片,因为这个可以获取user表里的所有数据,所以需要切片接收数据
	//db.Find(&user)      //获取user表里所有数据
	//fmt.Println(user)   //打印此时的user数据,是user表的所有数据
	////根据user表去获取CreditCard表的数据
	//db.Preload("CreditCards").Find(&user)
	//fmt.Println(user) //打印此时的user数据,是包含user表的数据和creditcard表的数据的
	////循环user数据,这一层循环能获取到user表的所有字段信息
	//for _, b := range user {
	//	fmt.Println(b.UserName)
	//	//根据user字段的信息去获取CreditCard表所有字段的信息
	//	for _, i := range b.CreditCards {
	//		fmt.Println(i.Number, b.UserName)
	//	}
	//}

	//第二种是Related, 这种方法在"github.com/jinzhu/gorm"才存在
	//var customer Customer //Related只能获取一条数据,所以不需要用切片接收
	//db.Find(&customer)    //只能获取user表的最后一条数据,赋值给user
	//fmt.Println(customer)
	////根据user表的信息,去获取CreditCard表的数据,由于只获取到的一条user表的信息,所以只能获取这条user表信息对应的CreditCard表的信息,赋值给user
	//db.Model(&customer).Related(&customer.CreditCards, "CreditCard").Find(&user)
	//fmt.Println(customer)
	//fmt.Println(customer.CreditCards)
	////根据这个user表的信息去获取对应CreditCard的所有字段
	//for _, i := range customer.CreditCards {
	//	fmt.Println(i.Number)
	//}
}
