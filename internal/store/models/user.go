package models

import "github.com/mattn/go-nulltype"

type User struct {
	ID        uint               `json:"-" gorm:"primaryKey" binding:"required"`
	SSOUserID nulltype.NullInt64 `json:"sso_user_id" gorm:"unique;not null" binding:"required"`
}
