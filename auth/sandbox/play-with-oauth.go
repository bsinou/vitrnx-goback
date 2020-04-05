package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/bsinou/vitrnx-goback/auth"
	"github.com/bsinou/vitrnx-goback/conf"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

const jwksWellKnownSuffix = ".well-known/jwks.json"

// var authServerURL = "http://localhost:8080/oidc/"
var jwks *jose.JSONWebKeySet

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return json.NewDecoder(strings.NewReader(bodyString)).Decode(target)
}

func getJWKS() (*jose.JSONWebKeySet, error) {
	if jwks != nil {
		return jwks, nil
	}

	jwks = &jose.JSONWebKeySet{}
	if err := getJSON(conf.GetJWKSURL(), jwks); err != nil {
		return nil, err
	}

	return jwks, nil
}

func VerifyToken(tokenStr string) (*jwt.Claims, error) {

	var claims jwt.Claims
	var token *jwt.JSONWebToken

	jwks, err := getJWKS()
	if err != nil {
		fmt.Println("could *not* retrieve public key:", err)
		return &claims, err
	}

	token, err = jwt.ParseSigned(tokenStr)
	if err != nil {
		fmt.Println("could not parse ID token", err)
		return &claims, err
	}

	fmt.Printf("Parsed signed token: %v\nHeaders: %v\n", token, token.Headers)

	err = token.Claims(jwks, &claims)
	if err != nil {
		fmt.Printf("could *not* decode claims, cause: %s\n", err.Error())
		return &claims, err
	}

	fmt.Printf("Payload?: Subject: %s, ID: %s, Expiry: %v\n", claims.Subject, claims.ID, claims.Expiry)
	auth.DebugJWTClaims(&claims)

	err = claims.Validate(jwt.Expected{
		Issuer: conf.GetAuthURL(),
		Time:   time.Now(),
	})
	if err != nil {
		fmt.Println("could *not* validate claims", err)
		return &claims, err
	}

	var claims2 map[string]interface{}
	err = token.Claims(jwks, &claims2)
	if err != nil {
		fmt.Println("could *not* decode custom claims", err)
		return &claims, err
	}

	fmt.Println("Custom claims:")
	for k, v := range claims2 {
		fmt.Printf("- %s: %v\n", k, v)
	}

	return &claims, err

}

func main() {

	fmt.Println("Testing valid token")
	var ok = "eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo5MDkyZGY4NC1jODY3LTRkYTAtODhmYi00MzhkYjQ0YTljMTQiLCJ0eXAiOiJKV1QifQ.eyJhdF9oYXNoIjoiX1BhYVNBcy1wOHo0dFlXZDE0bUNLQSIsImF1ZCI6WyJ2aXRybngiXSwiYXV0aF90aW1lIjoxNTc1NjMyNjUzLCJleHAiOjE1NzU2MzYyNTYsImlhdCI6MTU3NTYzMjY1NiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL29pZGMvIiwianRpIjoiZDc1YjlmMjItMWJhMy00MGJlLWE5ZGQtYzYyZjZhYzRmNjkxIiwibmFtZSI6ImFkbWluIiwibm9uY2UiOiIiLCJyYXQiOjE1NzU2MzI2NDMsInNpZCI6IjJhNzM1ODI4LWExY2ItNDAxMC05ZmY1LWNjYjcwNGU4NDFmOCIsInN1YiI6IjE2YWJjYjU2LWRmNTYtNGZiMy1hODM5LWU0NWU2ZjI0MGIxNyJ9.q1txbJB2bxhclguw14Gd7iWdgYCtQ6qksElq2BnP6v2jgo0KOLpgDhOViGFtcdAK5a26xpLu3TYkfQwd1K2LW9-X7whqCZ2ijo2O_Ho6_rq6sinlmlDy8ymEMNfPF_n--j9jXyOJK_hLwfCt05ySY21WpiqMvQSg6A_TSfPV0UWP7Al8uMasKR1z6wbaP1kzIJqLeFuYehB4R7UFYy6XU9ZZTRCi0BN1ENGrIGckj_v1e1qMpH_fuwisU0nSPOrSxkz-ov_xbwWiYA3vp3pjs8sF8RBDmP2ZswKZX1HQiJ1Gy0GBOkd6_fbCKS_bm3cpEuoH80y1PfJ0w9D-4crJ_pbWX48llZM8rBWvUF4I00hn0LoXETOzbbidZcb_cSVLIHrvjxw0-KWLZD6W_40ppewSe2spJRwDXP6Mgnlh7TMuKkCDwKGLvJS4WSp2FHxGI0IjHwi7TIt7PBcsiXwn9q9zjec-hRba4Ab8bXpwbWu-G_xoJFJq1cAyUpwsQ-O3EJMD_9v9zPX9GXF2unqf_8Au6OmJ31ebtIHhO7wui13X8VisGiGGlqw0DAQ2Demj1BKAeZVr7QZnkbfeokPt4unrsUIkw0wwIxZTIb8T5I1CpOaueJ5y3CAwFVK__JSvS_PDUhQ9MMcVvfBGEvMkj6d6xSVLmwOkBM8JIWi265A"
	_, err := VerifyToken(ok)
	if err != nil {
		fmt.Println("Test failed:", err)
	}

	fmt.Println("Testing *not* valid token")
	var nok = "eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzpkOWExNTVmNS02OTI1LTQ2MDktYjg0OS1mZTliOTdiNmU4MzUiLCJ0eXAiOiJKV1QifQ.eyJhdF9oYYNoIjoiU0IyQnhzQ01sZURFTkc1Y21la1FTdyIsImF1ZCI6WyJ2aXRybngiXSwiYXV0aF90aW1lIjoxNTc1MzAwNTYyLCJleHAiOjE1NzUzMDQxNzIsImlhdCI6MTU3NTMwMDU3MiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL29pZGMvIiwianRpIjoiNThlMGE3MWYtOTliOS00MjM5LWFlZjMtNGY5Y2E2YWI5ZjNhIiwibmFtZSI6ImFkbWluIiwibm9uY2UiOiIiLCJyYXQiOjE1NzUzMDA1NTUsInNpZCI6IjE1OTQ1YjQ2LTY3ZTktNDQ2NS05MmI3LTJhMmIxNTQ2ZDUwNSIsInN1YiI6IjFkNTU0ZjMyLWMwOTYtNGFkNC1hNjU2LTIwMzNhMmY5ZWUwNCJ9.XXvCbW7OPO1NmVgU0BA_Lf72tmBDn_eQfrhzz0XMY6ccMrCvohRgbA0g7jgA5A0lqOzGuem-xqXBYx3LchYr9Jz4fkIno6_B8i21S_2crE_NTwOmNkwkFVYvZeyGan16aFkg3R6eSGuhb9ihoo9vs8E_VYOZu7hcIyQA-pQafeR0MBlRUDqNx6mBNnZhyrZWNL5vWkUyZfEoEthaVDufPgdBzNCGUnK9xi4yxioRmV9CplI-eKaCaMTUNOwtbDHHohPTygC3Nbgs_I3A7nlB2cJRTJhXGhbhsTSnRi2AIHlMBrk_IXpNfVv5gxkM-IRhExgbY9Pt26yY1rENcw_-sMdNEqr7rQ8qTtP_KnUdf3Ckj9C_9GJxoA14UHeEKsSSaUB994VDTM-hhg6wCBbKASB044y0L4LPTzY2WzPz-54S5e_MIdGSwg2eXP0Ifayr4hCznT77zs_dGBK8J-GmdvOxa_Tv-ggYcH7DGAqcrMHiJJz6omptmvNvBn2A2ZqA7w_3KMnyWhai2uyogfUGpsV0rfYiVbPcPCTTOX-7fPcTsD_K1nZySn0ogo1421ECiEi5l2Yzt8eIaGF4FtdDMLrqtgcbv4xiS3W_ziAGoyyiYJtsc1HuZ1koE-QKU9ci0yWui9c08bwVTfMVSOAv6seFr2ctQ_OGMS5sfaHE-Gw"
	_, err = VerifyToken(nok)
	if err == nil {
		fmt.Println("Test failed: no error")
	}

}

// func main() {

// 	// jsonWebKey := jose.JSONWebKey{}
// 	// jsonWebKey.UnmarshalJSON([]byte(`
// 	// {
// 	// 	"use": "sig",
// 	// 	"kty": "RSA",
// 	// 	"kid": "public:d9a155f5-6925-4609-b849-fe9b97b6e835",
// 	// 	"alg": "RS256",
// 	// 	"n": "19iVeQNQ328MjG8kbAatRQcvhUVNj4tErvs4yt2Eqg_VqA79SI9aRVLJ_nbwNRLY4OII3kwEKklDg5AB_wfuI2rph-BFzWSzqjQkB37KsXYSdX7MWqw6wYrXV8ViPktqzRC1ZqOO-cPnodiNJV4IKDSzXt30U8vcxVDHZDBafQ68paaRslobktEegMQSiX0rzqogstvcrd1zRyUSYLZ8T23180U2TrvFeVRs2JwSbD7FsHGu7uH1TQGh0CjRph2H8NOGhXk1dVMD19OZJUYByGeTxWfGvmWnX48sTGsia4vfplKKdVuxKEk9otSyfc7tiN1SgR_ZntlzgmgkMX7qPlDiVcVgj5od7SCb7TrR6xxPjm7Fg_90KV7WS8JHhUztCIspHXu__w3Oskt8ALe62JW6oPFDduFsyVjYESFxn2Z-iWbJmSVZxnwdyFbP_7XLawjvrGlElLLrXpjuGzR2JEVyoom-n9S985qFUeThdHOX00Iyzcxf9hbAUFI4exNe82ALjpe3ohKMqVpEMnCi3xdi1A7UVOs3ciSMWiHqlK7pAhkhWkUpggD1nwGgELQ-ofa909EcogEc8tnBD9NvXhOIeWQzWgcvAB1Sk-yk6fiK-ypDVoBrYz5hCAEG0MahPuT_OIHa80us3W2mmhB5-_ILn2tKMkB6qFDun6XsGCc",
// 	// 	"e": "AQAB"
// 	// }`))

// 	jwks, err := getJWKS()
// 	if err != nil {
// 		fmt.Println("could *not* retrieve public key:", err)
// 	}

// 	// decode JWT token and verify signature using JSON Web Keyset
// 	fmt.Println("#### First test: OK")

// 	t3, err := jwt.ParseSigned(ok)
// 	if err != nil {
// 		fmt.Println("could not verify", err)
// 	} else {
// 		fmt.Println(" - ParseSigned: OK")
// 	}

// 	var claims map[string]interface{}
// 	err = t3.Claims(jwks, &claims)
// 	if err != nil {
// 		fmt.Println("could *not* decode claims", err)
// 	}

// 	fmt.Printf("Verified: %v\n", claims)

// 	fmt.Println("#### Second test: tempered")
// 	var claims2 map[string]interface{}

// 	t4, err := jwt.ParseSigned(nok)
// 	if err != nil {
// 		fmt.Println("could not verify", err)
// 		fmt.Println("#2: OK")
// 	} else {
// 		fmt.Printf("ERROR: parse signed passed on a tempered jwt: %v\n", t4)
// 	}

// 	err = t4.Claims(jwks, &claims2)
// 	if err != nil {
// 		fmt.Println("EXPECTED BEHAVIOUR: could *not* decode claims", err)
// 	}

// 	fmt.Printf("Second payload: %v\n", claims2)
// }
