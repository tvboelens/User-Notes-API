package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	Title  string `gorm:"not null"`
	Body   string
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
