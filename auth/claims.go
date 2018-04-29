// Package auth centralises the logic to add authentification mechanisms to the backend
package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/model"
)

// FIXME hard-coded known addresses
var knownAddresses = []string{
	"bruno@sinou.org",
	"irene@sinou.org",
	"pierre.cogne@gmail.com",
}

func WithClaims(ctx *gin.Context) {

	currUserName := ctx.MustGet(model.KeyUserName).(string)

	isAdmin := false
	for _, val := range knownAddresses {
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

	fmt.Println("############# Claims set")
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
