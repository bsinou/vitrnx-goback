// Package auth centralises the logic to add authentification mechanisms to the backend
package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/model"
)

func GetUserMeta(ctx *gin.Context) interface{} {
	// TODO implement
	return map[string]interface{}{
		// "email":       "bruno@sinou.org",
		// "displayName": "Bruno",
		// "roles":       []string{"ADMIN", "USER_ADMIN", "MODERATOR", "EDITOR", "VOLUNTEER", "GUEST"},
		"email":       "guest@sinou.org",
		"displayName": "Guest",
		"roles":       []string{"GUEST"},
	}
}

func WithClaims(ctx *gin.Context) {

	currUserName := ctx.MustGet(model.KeyUserName).(string)

	isAdmin := false
	for _, val := range viper.GetStringSlice(conf.KeyAdminUsers) {

		if currUserName == val {
			isAdmin = true
			break
		}
	}
	if isAdmin {
		ctx.Set(model.KeyClaims, []string{
			model.PolicyCanRead,
			model.PolicyCanEdit,
			model.PolicyCanManage,
		})
	} else {
		ctx.Set(model.KeyClaims, []string{
			model.PolicyCanRead,
		})
	}
}

// GetClaims returns an array with current valid claims to be serialised in JSON
func GetClaims(ctx *gin.Context) map[string]string {
	// BOAAAF enhance :)
	claims := ctx.MustGet(model.KeyClaims).([]string)
	claimMap := map[string]string{
		"canRead":   "false",
		"canEdit":   "false",
		"canManage": "false",
	}
	for _, claim := range claims {
		switch claim {
		case model.PolicyCanRead:
			claimMap["canRead"] = "true"
			break
		case model.PolicyCanEdit:
			claimMap["canEdit"] = "true"
			break
		case model.PolicyCanManage:
			claimMap["canManage"] = "true"
			break
		default:
			fmt.Println("Unknown claim: " + claim)
		}
	}
	fmt.Println("###### Getting claims for JSON")

	for k, v := range claimMap {
		fmt.Printf("%s : %s \n", k, v)

	}

	return claimMap
}
