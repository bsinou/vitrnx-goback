package gorm

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/model"

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

	sqliteDbPath := "./data/gorm-sqlite.db"
	if conf.Env != conf.EnvDev {
		sqliteDbPath = fmt.Sprintf("/var/lib/%s/data/gorm-sqlite.db", conf.VitrnxInstanceKey)
	}

	db, err := gorm.Open("sqlite3", sqliteDbPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot open sqlite at %s, : %s\n", sqliteDbPath, err))
	}
	db.LogMode(true) // Display SQL queries

	fmt.Printf("Initialised SQLite DB with file at %s\n", sqliteDbPath)
	return db
}

func createTableIfNeeded(db *gorm.DB) {
	if !db.HasTable(&model.User{}) {
		db.CreateTable(&model.User{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.User{})
	}
}
