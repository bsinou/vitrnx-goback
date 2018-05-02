package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
	"github.com/bsinou/vitrnx-goback/test/dummydata"
	"github.com/gin-gonic/gin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPostDao_CRUD(t *testing.T) {

	// insure the DB collection is empty before launching the tests
	cleanDB(t)

	Convey("TestBench", t, func() {
		ctx, _ := mockContext(dummydata.TestPost1)

		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)
		ctx.Set(model.KeyDb, db)

		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		posts := []model.Post{}
		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 1)

		cleanDB(t)

		posts2 := []model.Post{}
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
		So(err, ShouldBeNil)
		So(len(posts2), ShouldEqual, 0)

	})

	Convey("MultipleInserts", t, func() {
		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)

		ctx, _ := mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		ctx2, _ := mockContext(dummydata.TestPost2)
		ctx2.Set(model.KeyDb, db)
		PutPost(ctx2)
		So(len(ctx.Errors), ShouldEqual, 0)

		posts := []model.Post{}
		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 2)

		cleanDB(t)

		posts2 := []model.Post{}
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
		So(err, ShouldBeNil)
		So(len(posts2), ShouldEqual, 0)
	})

	Convey("Update", t, func() {
		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)

		ctx, _ := mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		// Update description
		{
			posts := []model.Post{}
			err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
			So(err, ShouldBeNil)
			So(len(posts), ShouldEqual, 1)

			updatePost := posts[0]
			newDesc := "A slightly different description"
			updatePost.Desc = newDesc
			b, err := json.Marshal(updatePost)
			So(err, ShouldBeNil)

			ctx, _ := mockContext(string(b))
			ctx.Set(model.KeyDb, db)
			PutPost(ctx)
			So(len(ctx.Errors), ShouldEqual, 0)

			posts = []model.Post{}
			err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
			So(err, ShouldBeNil)
			So(len(posts), ShouldEqual, 1)
			So(posts[0].Desc, ShouldEqual, newDesc)

		}

		// Try to update path => should fail
		{
			posts := []model.Post{}
			err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
			So(err, ShouldBeNil)
			So(len(posts), ShouldEqual, 1)

			updatePost := posts[0]
			newSlug := "a-slightly-different-slug"
			updatePost.Path = newSlug
			b, err := json.Marshal(updatePost)
			So(err, ShouldBeNil)

			ctx, _ := mockContext(string(b))
			ctx.Set(model.KeyDb, db)
			PutPost(ctx)
			So(len(ctx.Errors), ShouldEqual, 1)

			posts = []model.Post{}
			err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
			So(err, ShouldBeNil)
			So(len(posts), ShouldEqual, 1)
			So(posts[0].Path, ShouldEqual, "simple-test")
		}

		cleanDB(t)
	})

	Convey("TestPath", t, func() {

		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)

		ctx, _ := mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		// try to add same again
		ctx, _ = mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 1)

		posts := []model.Post{}
		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 1)

		exist := doesPathExist("simple-test", db)
		So(exist, ShouldBeTrue)

		exist = doesPathExist("simple-test-67", db)
		So(exist, ShouldBeFalse)

		cleanDB(t)
	})

	Convey("Queries", t, func() {
		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)

		ctx, _ := mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		ctx2, _ := mockContext(dummydata.TestPost2)
		ctx2.Set(model.KeyDb, db)
		PutPost(ctx2)
		So(len(ctx.Errors), ShouldEqual, 0)

		posts := []model.Post{}
		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 2)

		cleanDB(t)

		posts2 := []model.Post{}
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
		So(err, ShouldBeNil)
		So(len(posts2), ShouldEqual, 0)
	})

	Convey("TestDates", t, func() {

		s := dummydata.Session.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)

		ctx, _ := mockContext(dummydata.TestPost1)
		ctx.Set(model.KeyDb, db)
		t1 := time.Now()
		time.Sleep(20 * time.Millisecond)
		PutPost(ctx)
		time.Sleep(20 * time.Millisecond)
		So(len(ctx.Errors), ShouldEqual, 0)
		t2 := time.Now()

		posts := []model.Post{}
		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 1)

		fp := posts[0]

		So(fp.Date.After(t1), ShouldBeTrue)
		So(fp.Date.Before(t2), ShouldBeTrue)

		cleanDB(t)
	})

}

func mockContext(bodyAsJSONString string) (*gin.Context, *httptest.ResponseRecorder) {

	// See https://github.com/gin-gonic/gin/issues/323
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Set user mail that is normally added by a wrapper
	ctx.Set(model.KeyUserName, "test@example.com")
	ctx.Set(model.KeyUserId, "a very credible user ID")

	// Add the JSON body
	reader := bytes.NewReader([]byte(bodyAsJSONString))
	req, err := http.NewRequest(http.MethodPost, "/posts", ioutil.NopCloser(reader))
	if err != nil {
		fmt.Println("cannot create new mock request")
	}
	ctx.Request = req
	req.Header.Set("Content-Type", "application/json")

	return ctx, w
}

func cleanDB(t *testing.T) {
	s := dummydata.Session.Clone()
	defer s.Close()
	db := s.DB(dummydata.Mongo.Database)

	posts := []model.Post{}
	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	if err != nil {
		t.Error(err)
		return
	}

	for _, post := range posts {
		fmt.Printf("PostID: %s\n", post.ID.Hex())

		query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}
		err := db.C(model.PostCollection).Remove(query)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
