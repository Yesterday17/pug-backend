package models

type PipeConstructed struct {
	ModelIONDP

	Module string `json:"module"`
	Pipe   string `json:"pipe"`

	Arguments string `json:"arguments"`
}
