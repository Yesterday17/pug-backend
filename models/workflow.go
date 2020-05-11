package models

type Workflow struct {
	ModelIONDP

	Pipelines []Pipeline `json:"pipelines" gorm:"foreignKey:id"`

	ActiveWork  Work   `json:"active_work" gorm:"foreignKey:id"`
	HistoryWork []Work `json:"history_work" gorm:"foreignKey:id"`
}
