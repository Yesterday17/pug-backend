package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type InitSettings struct {
	Debug bool `json:"debug"`
}

func InitModels(s InitSettings) *gorm.DB {
	db, err := gorm.Open("sqlite3", "pug.db")

	if err != nil {
		log.Panic(err)
	}

	if s.Debug {
		db.Debug()
	}
	return db
}
