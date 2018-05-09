package gorm

import (
	"github.com/bsinou/vitrnx-goback/model"

	// First test with SQLite
	_ "github.com/mattn/go-sqlite3"
)

// WithUserRoles retrieves the roles stored as object in gorm and stores
// a string array of roles ID to ease JSON serialisation and manipulation by the front
// User must have been preloaded with the underlying roles
func WithUserRoles(user *model.User) {

	// err := db.Preload("Roles").Where(&model.User{UserID: userID}).First(&user).Error
	// if err != nil {
	// 	return userMeta, err
	// }

	roles := make([]string, len(user.Roles))

	for i, role := range user.Roles {
		roles[i] = role.RoleID
	}
	user.UserRoles = roles
}
