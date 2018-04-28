package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bsinou/vitrnx-goback/model"
	"github.com/bsinou/vitrnx-goback/mongodb"
	"github.com/gin-gonic/gin"

	. "github.com/smartystreets/goconvey/convey"
)

// func TestJwtToken(t *testing.T) {

// 	// See https://github.com/gin-gonic/gin/issues/323
// 	w := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(w)
// 	// Create a dummy gin context

// 	s := mongodb.Session.Clone()
// 	defer s.Close()
// 	ctx.Set(model.KeyDb, s.DB(mongodb.Mongo.Database))
// 	ctx.Set(model.KeyUserName, "test@example.com")

// 	// retrieve DB and store in context

// 	// test put and list
// }

func mockContext(params map[string]string) (*gin.Context, *httptest.ResponseRecorder) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)

	var paramsSlice []gin.Param
	for key, value := range params {
		paramsSlice = append(paramsSlice, gin.Param{
			Key:   key,
			Value: value,
		})
	}

	reader := bytes.NewReader([]byte(testJSONMsg))
	req, err := http.NewRequest(http.MethodPost, "/posts", ioutil.NopCloser(reader))
	if err != nil {
		fmt.Println("cannot create new mock request")
	}
	req.Header.Set("Content-Type", "application/json")
	context.Request = req

	context.Params = paramsSlice
	return context, w
}

func TestPostDao_CRUD(t *testing.T) {

	Convey("CreatePost", t, func() {
		ctx, _ := mockContext(map[string]string{
			"isPublic": "false",
			"title":    "A title",
			"path":     "some-crazy-thing",
			"tags":     "tag",
			"hero":     "test.jpg",
			"thumb":    "test.jpg",
		})

		s := mongodb.Session.Clone()
		defer s.Close()
		ctx.Set(model.KeyDb, s.DB(mongodb.Mongo.Database))
		ctx.Set(model.KeyUserName, "test@example.com")

		// fmt.Println("Got a json body: " + cardJSON)

		PutPost(ctx)

		// cardJSON := responseRecorder.Body.String()
		// So(cardJSON, ShouldBeBlank)
		// 		So(responseRecorder.Code, ShouldEqual, http.StatusNotFound)
	})
}

// { "post":
// }
var testJSONMsg string = `{ "body" : "",
		"desc" : "",
		"hero" : "test.jpg",
		"isPublic" : false,
		"path" : "abcd",
		"tags" : "",
		"thumb" : "test.jpg",
		"title" : "Un autre titre"}`
