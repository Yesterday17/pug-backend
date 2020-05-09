package models

type ModuleRestrictRule struct {
	Model

	ModuleName            string `json:"module" gorm:"not null"`
	PipeName              string `json:"pipe"`
	MinAvailableUserLevel int    `json:"level" gorm:"not null"`
}
