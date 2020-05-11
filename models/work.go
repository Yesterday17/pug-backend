package models

type Work struct {
	ModelWithID
}

type WorkLog struct {
	ModelWithID

	Log string `json:"log"`
}
