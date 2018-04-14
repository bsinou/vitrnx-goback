// Package gateway centralises logic to connect to the outer world
package gateway

import (
	"context"
	"testing"
)

const (
	jwtToTest = "A valid JWT token"
)

func TestJwtToken(t *testing.T) {
	err := CheckCredentialAgainstFireBase(context.Background(), jwtToTest)
	if err != nil {
		t.Error("invalid JWT", err.Error())
	}
}
