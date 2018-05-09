package auth

import (
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	fbauth "firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	jgorm "github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"
)

var (
	cfPath         string
	adminEmail     = "Not a valid email, should be overwritten using conf"
	anonymousEmail = "Not a valid email, should be overwritten using conf"
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
	meta, err := GetUserMeta(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(201, gin.H{"userMeta": meta})
}

// CheckCredentialAgainstFireBase simply validate the passed token against firebase.
func CheckCredentialAgainstFireBase(ctx *gin.Context, jwt string) error { //, uid

	// credOption := option.WithCredentialsFile(credFilePath())
	// fbApp, err := firebase.NewApp(ctx, nil, credOption)
	// // TODO add retry
	// if err != nil {
	// 	return fmt.Errorf("cannot connect to firebase: %v", err)
	// }
	// client, err := fbApp.Auth(ctx)
	// if err != nil {
	// 	return fmt.Errorf("error getting Auth client: %v", err)
	// }
	client, err := getFireBaseClient(ctx)
	if err != nil {
		return err
	}

	token, err := client.VerifyIDToken(jwt)
	if err != nil {
		return fmt.Errorf("JWT validation failed:  %v", err)
	}

	// Store relevant user info in the context
	cs := token.Claims
	ctx.Set(model.KeyUserID, cs[model.FbKeyUserID].(string))
	ctx.Set(model.KeyEmailVerified, cs[model.FbKeyEmailVerified].(bool))

	return nil
}

// ListExistingUsers retrieves all users from firebase
func ListExistingUsers(ctx *gin.Context) error { //, uid
	client, err := getFireBaseClient(ctx)
	if err != nil {
		return err
	}

	userIterator := client.Users(ctx, "")
	if err != nil {
		return fmt.Errorf("could not list users: %v", err)
	}

	// This must be enhanced, many shortcuts and hacks here...
	db := gorm.GetConnection()
	defer db.Close()
	ae := viper.GetString(conf.KeyAdminEmail)
	if ae != "" {
		adminEmail = ae
	}
	an := viper.GetString(conf.KeyAnonymousEmail)
	if an != "" {
		anonymousEmail = an
	}

	for {
		userRecord, err := userIterator.Next()
		if err != nil {
			if err.Error() == "no more items in iterator" {
				fmt.Println("Sync with firebase done.")
				break
			}
			return err
		}

		err = updateUser(db, userRecord)
		if err != nil {
			return err
		}
	}

	if userIterator.PageInfo().Remaining() > 0 {
		// TODO implement this
		return fmt.Errorf("pagination is not implemented and user count "+
			"is greater than what can fit in a page (%d users), missed %d users",
			userIterator.PageInfo().MaxSize, userIterator.PageInfo().Remaining())
	}

	return nil
}

/* HELPER FUNCTIONS */

func getFireBaseClient(ctx *gin.Context) (*fbauth.Client, error) {
	credOption := option.WithCredentialsFile(credFilePath())
	fbApp, err := firebase.NewApp(ctx, nil, credOption)
	// TODO add retry
	if err != nil {
		return nil, fmt.Errorf("cannot connect to firebase: %v", err)
	}
	client, err := fbApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}
	return client, nil
}

// Update the user repository if necessary
// TODO clean this using channel
func updateUser(db *jgorm.DB, eur *fbauth.ExportedUserRecord) error {

	var user model.User
	err := db.Where(&model.User{UserID: eur.UID}).First(&user).Error
	if err == nil {
		// User already exist do nothing
		return nil
	}
	if !jgorm.IsRecordNotFoundError(err) {
		return err // Unexpected error
	}

	roleStr := "REGISTERED"
	if eur.Email == adminEmail { // Create the admin user
		roleStr = "ADMIN"
	} else if eur.Email == anonymousEmail { // Create the admin user
		roleStr = "ANONYMOUS"
	}

	user = model.User{
		UserID: eur.UID,
		Email:  eur.Email,
		Name:   eur.DisplayName,
		Roles: []model.Role{
			{RoleID: roleStr},
		},
	}

	// db.Set("gorm:association_autoupdate", false).Save(&user)
	db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Save(&user)

	return nil
}

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
