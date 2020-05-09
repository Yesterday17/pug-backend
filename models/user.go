package models

type User struct {
	Model

	UUID      string `json:"uuid" gorm:"primary_key"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	UserLevel int    `json:"level" gorm:"not null"`

	Setting UserSettings `json:"setting" gorm:"foreignKey:uuid"`
}

type UserAccountSettings struct {
	Model

	UUID  string `json:"-" gorm:"primary_key"`
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email"`
	Icon  string `json:"icon"`
}

type UserSettings struct {
	Model
	UUID string `json:"-" gorm:"primary_key"`

	Account UserAccountSettings `json:"account" gorm:"foreignKey:uuid"`
}
