package models

import (
	"time"

	"github.com/mattn/go-nulltype"
	"gorm.io/gorm"
)

type Match struct {
	ID        uint               `json:"-" gorm:"primaryKey" binding:"required"`
	SSOUserID nulltype.NullInt64 `json:"sso_user_id,omitempty" gorm:"uniqueIndex"`
	Admin     bool               `json:"admin" gorm:"default:false"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"-"`
	DeletedAt gorm.DeletedAt     `json:"-" gorm:"index"`
}
