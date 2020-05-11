package models

type PipeConstructed struct {
	ModelWithID

	Owner     string `json:"-" gorm:"not null"`
	OwnerUser User   `json:"owner" gorm:"foreignKey:uuid"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`

	Arguments map[string]interface{} `json:"arguments" sql:"type:LONGTEXT"`
}
