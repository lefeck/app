package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const (
	UserAssociation         = "Users"
	UserAuthInfoAssociation = "AuthInfos"
	GroupAssociation        = "Groups"
)

type BaseModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"` // soft delete
}

type User struct {
	ID         uint     `json:"id" gorm:"autoIncrement;primaryKey"`
	Name       string   `json:"name" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	Password   string   `json:"-" gorm:"type:varchar(256);"`
	RePassword string   `json:"-" gorm:"type:varchar(256);"`
	Email      string   `json:"email" gorm:"type:varchar(256);"`
	Avatar     string   `json:"avatar" gorm:"type:varchar(256);"` // 头像
	UserInfo   UserInfo `json:"authInfo" gorm:"foreignKey:UserID;references:ID"`
	BaseModel
}

func (*User) TableName() string {
	return "users"
}

func (u *User) CacheKey() string {
	return u.TableName() + ":id"
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

//// 输出精简信息
//func (user *User) ToMapSimple(token string) map[string]interface{} {
//	return map[string]interface{}{
//		"user_id":     user.ID,
//		"user_name":   user.Name,
//		"user_avator": user.UserInfo.UserAvator,
//		"token":       token,
//	}
//}

// 输出详细信息
//func (user *User) ToMap() map[string]interface{} {
//	return map[string]interface{}{
//		"user_name":   user.Name,
//		"nick_name":   user.UserInfo,
//		"user_id":     user.ID,
//		"is_active":   user.IsActive,
//		"user_avator": user.UserInfo.UserAvator,
//		"user_email":  user.UserInfo.UserEmail,
//		"user_desc":   user.UserInfo.UserDesc,
//		"user_addr":   user.UserInfo.UserAddr,
//	}
//}

//// 输出详细信息 带token
//func (user *User) ToMapHasToken(token string) map[string]interface{} {
//	resp := user.ToMap()
//	resp["token"] = token
//	return resp
//}
