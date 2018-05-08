package gorm

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/model"

	// First test with SQLite
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbFileAbsPath string
)

// InitGormRepo initialises a Gorm backend backed by SQLite
func InitGormRepo() {
	dataDirPath := conf.GetDataFolderPath()
	dbFileAbsPath = filepath.Join(dataDirPath, "gorm-sqlite.db")
	db := getDb()
	createTableIfNeeded(db)
}

// InitGormTestRepo is a convenience method to ease Unit Test implementation
func InitGormTestRepo(sqliteFileAbsPath string) {
	dbFileAbsPath = sqliteFileAbsPath
	db := getDb()
	createTableIfNeeded(db)
}

// GetConnection manage a pool of session TODO
func GetConnection() *gorm.DB {
	return getDb()
}

/* Local helpers */
func getDb() *gorm.DB {

	db, err := gorm.Open("sqlite3", dbFileAbsPath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot open sqlite at %s, : %s\n", dbFileAbsPath, err))
	}
	// db.LogMode(true) // Display SQL queries

	fmt.Printf("Initialised SQLite DB with file at %s\n", dbFileAbsPath)
	return db
}

func createTableIfNeeded(db *gorm.DB) {
	if !db.HasTable(&model.User{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
		db.CreateTable(&model.User{}, &model.Role{})
	}
}
