package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/bsinou/vitrnx-goback/model"
	"github.com/gin-gonic/gin"
)

func mockContext() (*gin.Context, *httptest.ResponseRecorder) {
	// See https://github.com/gin-gonic/gin/issues/323
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Set user meta that are normally added by a wrapper
	ctx.Set(model.KeyUserDisplayName, "test")
	ctx.Set(model.KeyUserID, "a very credible user ID")

	return ctx, w
}

func withDummyRequest(ctx *gin.Context, bodyAsJSONString string) {
	// Add the JSON body
	reader := bytes.NewReader([]byte(bodyAsJSONString))
	req, err := http.NewRequest(http.MethodPost, "/posts", ioutil.NopCloser(reader))
	if err != nil {
		fmt.Println("cannot create new mock request")
	}
	req.Header.Set("Content-Type", "application/json")

	ctx.Request = req
}

// Simplify implementation of users tests
func getPreparedUser(name, mainRole string) *model.User {

	return &model.User{
		Name:  name,
		Email: fmt.Sprintf("%s@example.com", name),
		Roles: []model.Role{
			{
				RoleID: mainRole,
				Label:  mainRole,
			},
		},
	}
}
