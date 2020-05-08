package models

type User struct {
	Model

	UUID      string `json:"uuid" gorm:"unique;not null"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	UserLevel int    `json:"user_level" gorm:"not null"`
}

type UserAccountSettings struct {
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email"`
	Icon  string `json:"icon"`
}

type UserSettings struct {
	Model

	UUID    string              `json:"-" gorm:"unique;not null"`
	Account UserAccountSettings `json:"account"`
}
