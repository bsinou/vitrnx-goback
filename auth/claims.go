// Package auth centralises the logic to add authentification mechanisms to the backend
package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/model"
)

func GetUserMeta(ctx *gin.Context) map[string]interface{} {
	// TODO implement
	return map[string]interface{}{
		model.KeyEmail:           "bruno@sinou.org",
		model.KeyUserDisplayName: "Bruno",
		model.KeyRoles:           []string{"ADMIN", "USER_ADMIN", "MODERATOR", "EDITOR", "VOLUNTEER", "GUEST"},
		// model.KeyEmail:           "guest@sinou.org",
		// model.KeyUserDisplayName: "Guest",
		// model.KeyRoles:           []string{"GUEST"},
	}
}

func WithUserMeta(ctx *gin.Context) {

	userMeta := GetUserMeta(ctx)

	ctx.Set(model.KeyEmail, userMeta[model.KeyEmail])
	ctx.Set(model.KeyUserDisplayName, userMeta[model.KeyUserDisplayName])
	ctx.Set(model.KeyRoles, userMeta[model.KeyRoles])
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
