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

type AuthInfo struct {
	ID           uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	UserId       uint      `json:"userId" gorm:"size:256"`
	Url          string    `json:"url" gorm:"size:256"`
	AuthType     string    `json:"authType" gorm:"size:256"`
	AuthId       string    `json:"authId" gorm:"size:256"`
	AccessToken  string    `json:"-" gorm:"size:256"`
	RefreshToken string    `json:"-" gorm:"size:256"`
	Expiry       time.Time `json:"-"`
	BaseModel
}

type BaseModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"` // soft delete
}

type User struct {
	ID         uint       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name       string     `json:"name" gorm:"type:varchar(50);uniqueIndex:name_auth_type;not null"`
	Password   string     `json:"-" gorm:"type:varchar(256);"`
	RePassword string     `json:"-" gorm:"type:varchar(256);"`
	Email      string     `json:"email" gorm:"type:varchar(256);"`
	AuthType   string     `json:"authType" gorm:"type:varchar(256);uniqueIndex:name_auth_type;default:nil"`
	AuthId     string     `json:"authId" gorm:"type:varchar(256);"`
	Avatar     string     `json:"avatar" gorm:"type:varchar(256);"` // images
	AuthInfos  []AuthInfo `json:"authInfos" gorm:"foreignKey:UserId;references:ID"`
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
