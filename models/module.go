package models

type ModuleRestrictRule struct {
	Model

	ModuleName            string `json:"module" gorm:"not null"`
	PipeName              string `json:"pipe"`
	MinAvailableUserLevel int    `json:"user_level" gorm:"not null"`
}
