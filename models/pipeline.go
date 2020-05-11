package models

type Pipeline struct {
	ModelWithID

	Owner     string `json:"-" gorm:"not null"`
	OwnerUser User   `json:"owner" gorm:"foreignKey:owner"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`

	Pipes []PipeConstructed `json:"pipes" gorm:"foreignKey:id"`
}
