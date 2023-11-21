package main

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB(ctx context.Context) *gorm.DB {
	return db
}

type User struct {
	ID   int64
	Name string
	Age  int64
}

func GetUserByUserID(ctx context.Context, userID int64) (*User, error) {
	db := GetDB(ctx)
	var user User
	if userID > 0 {
		db = db.Where(`userID = ?`, userID)
	}
	if err := db.Model(&User{}).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsersByName(ctx context.Context, name string) (*User, error) {
	db := GetDB(ctx)
	var users User
	if name != "" {
		db = db.Where(`name like '%%'`, name)
	}
	if err := db.Model(&User{}).Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func GetUsersByAge(ctx context.Context, age int64) (*User, error) {
	db := GetDB(ctx)
	var user User
	if age > 0 {
		db = db.Where(`age = ?`, age)
	}
	if err := db.Model(&User{}).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

type UserParam struct {
	ID   int64
	Name string
	Age  int64
}

func GetUserInfos(ctx context.Context, info *UserParam) ([]*User, error) {
	db := GetDB(ctx)
	db = db.Model(&User{})
	var infos []*User
	if info.ID > 0 {
		db = db.Where(`age = ?`, info.ID)
	}
	if info.Name != "" {
		db = db.Where(`name = ?`, info.Name)
	}
	if info.Age > 0 {
		db = db.Where(`age = ?`, info.Age)
	}

	if err := db.Find(&info).Error; err != nil {
		return nil, err
	}
	return infos, nil
}
