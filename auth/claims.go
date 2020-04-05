// Package auth centralises the logic to add authentification mechanisms to the backend.
package auth

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/model"
)

// GetUserMeta retrieves current user metadata.
func GetUserMeta(ctx *gin.Context) (map[string]interface{}, error) {
	log.Println("Retrieving user Meta...")

	db := ctx.MustGet(model.KeyDb).(*gorm.DB)
	userID := ctx.MustGet(model.KeyUserID).(string)

	userMeta := make(map[string]interface{})

	var user model.User
	err := db.Preload("Roles").Where(&model.User{UserID: userID}).First(&user).Error
	if err != nil {
		return userMeta, err
	}

	userMeta[model.KeyEmail] = user.Email
	userMeta[model.KeyUserDisplayName] = user.Name

	roles := make([]string, len(user.Roles))

	for i, role := range user.Roles {
		roles[i] = role.RoleID
	}
	userMeta[model.KeyUserRoles] = roles

	return userMeta, nil
}

// WithUserMeta simply enriches current context with user metadata.
func WithUserMeta(ctx *gin.Context) error {

	userMeta, err := GetUserMeta(ctx)
	if err != nil {
		return err
	}

	ctx.Set(model.KeyEmail, userMeta[model.KeyEmail])
	ctx.Set(model.KeyUserDisplayName, userMeta[model.KeyUserDisplayName])
	ctx.Set(model.KeyUserRoles, userMeta[model.KeyUserRoles])

	return nil
}
