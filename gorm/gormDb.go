package gorm

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

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

	return db
}

func createTableIfNeeded(db *gorm.DB) {

	if !db.HasTable(&model.User{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
		db.CreateTable(&model.User{}, &model.Role{})
	}

	// initialise roles from config
	knownRoles := viper.GetStringSlice(conf.KeyKnownRoles)
	// Simplify testing
	if len(knownRoles) == 0 {
		knownRoles = []string{"ADMIN/Administrator", "REGISTERED/Registered User", "Anonymous/Anonymous User"}
	}

	for _, v := range knownRoles {
		tokens := strings.Split(v, "/")
		var role model.Role
		err := db.Where(&model.Role{RoleID: tokens[0]}).First(&role).Error
		if err == nil {
			// Role already exist do nothing
			continue
		}
		if !gorm.IsRecordNotFoundError(err) {
			log.Fatalln(err.Error()) // Unexpected error
		}
		db.Save(&model.Role{RoleID: tokens[0], Label: tokens[1]})
	}

	fmt.Printf("Initialised SQLite DB with file at %s\n", dbFileAbsPath)

	// vknownRoles := viper.GetStringSlice(conf.KeyKnownRoles)
	// fmt.Printf("#####  vknownRoles length %d\n ", len(vknownRoles))
	// for i, v := range vknownRoles {
	// 	fmt.Printf("#%d: %s\n", i, v)
	// }
}
