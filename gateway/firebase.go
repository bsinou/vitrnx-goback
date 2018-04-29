// Package gateway centralises logic to connect to the outer world
package gateway

import (
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"github.com/bsinou/vitrnx-goback/model"
)

const (
	credFilePath = "/home/bsinou/dev/private/vitrnx/firebase-apiCert.json"
)

var (
	fbApp *firebase.App
)

func init() {
}

// CheckCredentialAgainstFireBase simply validate the passed token against firebase.
func CheckCredentialAgainstFireBase(ctx *gin.Context, jwt string) error { //, uid

	credOption := option.WithCredentialsFile(credFilePath)
	var err error

	fbApp, err = firebase.NewApp(ctx, nil, credOption)
	// TODO add retry
	if err != nil {
		log.Fatalf("cannot connect to firebase: %v\n", err)
	}

	client, err := fbApp.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return err
	}

	token, err := client.VerifyIDToken(jwt)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return err
	}

	// Store relevant user info in the context
	cs := token.Claims
	ctx.Set(model.KeyUserId, cs[model.FbKeyUserId].(string))
	// We use user email as name for the time being
	ctx.Set(model.KeyUserName, cs[model.FbKeyEmail].(string))
	ctx.Set(model.KeyEmailVerified, cs[model.FbKeyEmailVerified].(bool))

	setClaims(ctx)

	return nil
}

// FIXME hard-coded known addresses
var knownAddresses = []string{
	"bruno@sinou.org",
	"irene@sinou.org",
	"pierre.cogne@gmail.com",
}

func setClaims(ctx *gin.Context) {

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
