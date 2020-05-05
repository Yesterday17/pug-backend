package models

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ModelSettings struct {
	Debug bool

	DBType   string
	DBConfig map[string]string
}

func InitModels(s *ModelSettings) (db *gorm.DB) {
	var err error

	switch s.DBType {
	// TODO: Support more db types
	case "sqlite":
		db, err = gorm.Open("sqlite3", s.DBConfig["db_name"])
	default:
		err = errors.New("unsupported database type")
	}
	if err != nil {
		log.Fatal("Failed to load database", err)
	}

	if s.Debug {
		db.Debug()
	}

	db.AutoMigrate(&User{})
	// db.AutoMigrate(&UserSettings{})
	return
}
