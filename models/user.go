package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	UUID string `json:"uuid",gorm:"unique;not null"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type UserSettings struct {
	gorm.Model
}
