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
	cfPath string
)

func init() {
}

// // Login delegates authentication to firebase.
// TODO this is not obvious: firebase does not ease direct auth with login / pwd from the go SDK
// I presume this is to insure they can collect end-user info (like IPADDRESS and client browser...) upon login
// func Login(ctx *gin.Context) error {

// 	client := fbClient(ctx)
// 	token, err := client.VerifyIDToken(jwt)
// 	if err != nil {
// 		log.Printf("error verifying ID token: %v\n", err)
// 		return err
// 	}

// 	// Store relevant user info in the context
// 	cs := token.Claims
// 	ctx.Set(model.KeyUserID, cs[model.FbKeyUserID].(string))
// 	// We use user email as name for the time being
// 	ctx.Set(model.KeyUserName, cs[model.FbKeyEmail].(string))
// 	ctx.Set(model.KeyEmailVerified, cs[model.FbKeyEmailVerified].(bool))

// 	WithClaims(ctx)

// 	return nil
// }

// PostLogin add vitrnx specific user info upon login
func PostLogin(ctx *gin.Context) {
	// TODO implement this
	ctx.JSON(201, gin.H{"userMeta": GetUserMeta(ctx)})
}

// CheckCredentialAgainstFireBase simply validate the passed token against firebase.
func CheckCredentialAgainstFireBase(ctx *gin.Context, jwt string) error { //, uid

	credOption := option.WithCredentialsFile(credFilePath())
	fbApp, err := firebase.NewApp(ctx, nil, credOption)
	// TODO add retry
	if err != nil {
		log.Fatalf("cannot connect to firebase: %v\n", err)
	}
	client, err := fbApp.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(jwt)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return err
	}

	// Store relevant user info in the context
	cs := token.Claims
	ctx.Set(model.KeyUserID, cs[model.FbKeyUserID].(string))
	ctx.Set(model.KeyEmailVerified, cs[model.FbKeyEmailVerified].(bool))
	WithUserMeta(ctx)
	return nil
}

/* HELPER FUNCTIONS */
// func fbClient(ctx *gin.Context) *auth.Client {

// 	return client
// }

// Caches path to local Firebase API cert file
func credFilePath() string {

	if cfPath != "" {
		return cfPath
	}

	var err error
	cfPath, err = conf.GetConfigFile("firebase-apiCert.json")
	if err != nil {
		log.Fatalf("no firebase API cert file found")
	}
	return cfPath
}
