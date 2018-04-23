// Package gateway centralises logic to connect to the outer world
package gateway

import (
	"context"
	"log"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

const (
	credFilePath = "/home/bsinou/dev/private/vitrnx/firebase-apiCert.json"
)

var (
	fbApp *firebase.App
)

func init() {
	credOption := option.WithCredentialsFile(credFilePath)
	var err error
	fbApp, err = firebase.NewApp(context.Background(), nil, credOption)
	// TODO add retry
	if err != nil {
		log.Fatalf("cannot connect to firebase: %v\n", err)
	}
}

// CheckCredentialAgainstFireBase simply validate the passed token against firebase.
func CheckCredentialAgainstFireBase(ctx context.Context, jwt string) error {

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

	log.Printf("Verified ID token: %v\n", token)

	return nil
}
