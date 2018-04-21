package gorm

import (
	"github.com/bsinou/vitrnx-goback/model"
	"github.com/jinzhu/gorm"

	// First test with SQLite
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	db := getDb()
	createTableIfNeeded(db)
}

// GetConnection manage a pool of session TODO
func GetConnection() *gorm.DB {
	return getDb()
}

/* Local helpers */

func getDb() *gorm.DB {
	// DB: launch and config
	db, err := gorm.Open("sqlite3", "./data/gorm-sqlite.db")
	if err != nil {
		panic(err) // TODO enhance
	}
	db.LogMode(true) // Display SQL queries
	return db
}

func createTableIfNeeded(db *gorm.DB) {
	if !db.HasTable(&model.User{}) {
		db.CreateTable(&model.User{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.User{})
	}
}
