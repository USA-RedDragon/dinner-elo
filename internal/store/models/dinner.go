package models

type Dinner struct {
	ID   uint   `json:"-" gorm:"primaryKey" binding:"required"`
	Name string `json:"name" gorm:"uniqueIndex" binding:"required"`
}
