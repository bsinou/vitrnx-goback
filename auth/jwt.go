package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/bsinou/vitrnx-goback/conf"
)

var (
	knownAdmins []string
	jwks        *jose.JSONWebKeySet
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func init() {
}

type CustomClaims struct {
	*jwt.Claims
	Login       string
	DisplayName string
	Email       string
	Roles       []string
}

func VerifyToken(tokenStr string) (*CustomClaims, error) {

	var claims jwt.Claims
	var cclaims CustomClaims
	var token *jwt.JSONWebToken

	jwks, err := getJWKS()
	if err != nil {
		fmt.Println("could *not* retrieve public key:", err)
		return &cclaims, err
	}

	token, err = jwt.ParseSigned(tokenStr)
	if err != nil {
		fmt.Println("could not parse ID token", err)
		return &cclaims, err
	}

	fmt.Printf("Parsed signed token: %v\nHeaders: %v\n", token, token.Headers)

	err = token.Claims(jwks, &claims)
	if err != nil {
		fmt.Println("could *not* decode claims", err)
		return &cclaims, err
	}

	// The issuer claim generally ends with a "/"
	iss := conf.GetAuthURL() + "/"
	err = claims.Validate(jwt.Expected{
		Issuer: iss,
		Time:   time.Now(),
	})
	if err != nil {
		fmt.Println("validation failed", err)
		return &cclaims, err
	}

	cclaims.Claims = &claims

	DebugJWTClaims(&claims)
	// TODO enhance this
	var claims2 map[string]interface{}
	err = token.Claims(jwks, &claims2)
	if err != nil {
		fmt.Println("could *not* decode custom claims", err)
		return &cclaims, err
	}
	cclaims.Login = claims2["name"].(string)
	cclaims.DisplayName = claims2["name"].(string)

	roles := []string{"ANONYMOUS"}

	if HasRole(conf.GetKnownAdmins(), cclaims.Login) {
		roles = []string{"ADMIN"}
	}
	cclaims.Roles = roles

	return &cclaims, err
}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func getJWKS() (*jose.JSONWebKeySet, error) {
	if jwks != nil {
		return jwks, nil
	}

	jwks = &jose.JSONWebKeySet{}
	if err := getJSON(conf.GetJWKSURL(), jwks); err != nil {
		fmt.Printf("could not retrieve JWKS from %s, cause: %s \n", conf.GetJWKSURL(), err.Error())
		fmt.Printf("##### Most probably, your authURL is not correctly configured, check your config files.\n")
		return nil, err
	}

	return jwks, nil
}

func HasRole(roles []string, toCheck string) bool {
	for _, a := range roles {
		if a == toCheck {
			return true
		}
	}
	return false
}

const debugFormat = "2006/01/02 - 15:04:05 MST"

func DebugJWTClaims(c *jwt.Claims) {

	fmt.Println("Debuging claims:")
	fmt.Printf("Issuer: %s\n", c.Issuer)
	fmt.Printf("Subject: %s\n", c.Subject)
	fmt.Printf("Audience: %v\n", c.Audience)
	fmt.Printf("Expiry: %s\n", c.Expiry.Time().Format(debugFormat))
	fmt.Printf("NotBefore: %s\n", c.NotBefore.Time().Format(debugFormat))
	fmt.Printf("IssuedAt: %s\n", c.IssuedAt.Time().Format(debugFormat))
	fmt.Printf("ID: %s\n", c.ID)
}
