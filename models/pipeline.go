package models

type Pipeline struct {
	ModelIONDP

	Pipes []PipeConstructed `json:"pipes" gorm:"foreignKey:id"`
}
