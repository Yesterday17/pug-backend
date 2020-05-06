package models

type User struct {
	Model

	UUID     string `json:"uuid" gorm:"unique;not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
	Icon     string `json:"icon"`
}

type UserSettings struct {
	Model
}
