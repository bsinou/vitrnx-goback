// Package auth centralises the logic to add authentification mechanisms to the backend
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/model"
)

func GetUserMeta(ctx *gin.Context) (map[string]interface{}, error) {
	userID := ctx.MustGet(model.KeyUserID).(string)
	db := ctx.MustGet(model.KeyUserDb).(*gorm.DB)

	var userMeta map[string]interface{}

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

// func WithClaims(ctx *gin.Context) {

// 	currId := ctx.MustGet(model.KeyUserID).(string)

// 	isAdmin := false
// 	for _, val := range viper.GetStringSlice(conf.KeyAdminUsers) {

// 		if currUserName == val {
// 			isAdmin = true
// 			break
// 		}
// 	}
// 	if isAdmin {
// 		ctx.Set(model.KeyClaims, []string{
// 			model.PolicyCanRead,
// 			model.PolicyCanEdit,
// 			model.PolicyCanManage,
// 		})
// 	} else {
// 		ctx.Set(model.KeyClaims, []string{
// 			model.PolicyCanRead,
// 		})
// 	}
// }

// // GetClaims returns an array with current valid claims to be serialised in JSON
// func GetClaims(ctx *gin.Context) map[string]string {
// 	// BOAAAF enhance :)
// 	claims := ctx.MustGet(model.KeyClaims).([]string)
// 	claimMap := map[string]string{
// 		"canRead":   "false",
// 		"canEdit":   "false",
// 		"canManage": "false",
// 	}
// 	for _, claim := range claims {
// 		switch claim {
// 		case model.PolicyCanRead:
// 			claimMap["canRead"] = "true"
// 			break
// 		case model.PolicyCanEdit:
// 			claimMap["canEdit"] = "true"
// 			break
// 		case model.PolicyCanManage:
// 			claimMap["canManage"] = "true"
// 			break
// 		default:
// 			fmt.Println("Unknown claim: " + claim)
// 		}
// 	}
// 	fmt.Println("###### Getting claims for JSON")

// 	for k, v := range claimMap {
// 		fmt.Printf("%s : %s \n", k, v)

// 	}

// 	return claimMap
// }
