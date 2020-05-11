package models

type PipeConstructed struct {
	ModelWithID

	Owner     string `json:"-" gorm:"not null"`
	OwnerUser User   `json:"owner" gorm:"foreignKey:owner"`

	Module string `json:"module"`
	Pipe   string `json:"pipe"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`

	Arguments string `json:"arguments"`
}
