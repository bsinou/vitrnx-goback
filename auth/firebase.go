package auth

import (
	"log"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/model"
)

var (
	fbApp *firebase.App
)

func init() {
}

// CheckCredentialAgainstFireBase simply validate the passed token against firebase.
func CheckCredentialAgainstFireBase(ctx *gin.Context, jwt string) error { //, uid

	credFilePath, err := conf.GetConfigFile("firebase-apiCert.json")
	if err != nil {
		log.Fatalf("no firebase API cert file found")
	}
	credOption := option.WithCredentialsFile(credFilePath)

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

	WithClaims(ctx)

	return nil
}