package model

import "time"

type UserInfo struct {
	ID           uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	UserID       uint      `json:"userID" gorm:"size:256"`
	AccessToken  string    `json:"-" gorm:"size:256"`
	RefreshToken string    `json:"-" gorm:"size:256"`
	Expiry       time.Time `json:"-"`
	BaseModel
}
