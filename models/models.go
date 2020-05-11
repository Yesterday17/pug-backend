package models

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ModelSettings struct {
	Debug bool `json:"debug"`

	DBType   string            `json:"db_type"`
	DBConfig map[string]string `json:"db_config"`
}

type Model struct {
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

type ModelWithHiddenID struct {
	Model

	ID uint `json:"-" gorm:"primary_key"`
}

type ModelWithID struct {
	Model

	ID uint `json:"id" gorm:"primary_key"`
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
	db.AutoMigrate(&UserSettings{})
	db.AutoMigrate(&UserAccountSettings{})

	db.AutoMigrate(&ModuleRestrictRule{})

	db.AutoMigrate(&PipeConstructed{})
	return
}
