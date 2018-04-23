// Package gateway centralises logic to connect to the outer world
package gateway

import (
	"context"
	"testing"
)

const (
	jwtToTest = ""
)

func TestJwtToken(t *testing.T) {

	if jwtToTest != "" {
		err := CheckCredentialAgainstFireBase(context.Background(), jwtToTest)
		if err != nil {
			t.Error("invalid JWT", err.Error())
		}
	}
}
