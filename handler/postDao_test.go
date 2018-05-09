package handler

import (
	"encoding/json"
	"testing"
	// "time"

	// "gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
	"github.com/bsinou/vitrnx-goback/test/dummydata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostDao_CRUD(t *testing.T) {

	// insure the DB collection is empty before launching the tests
	cleanMongoDB(t)

	Convey("Simple test", t, func() {
		ctx, _ := mockContext()

		s := dummydata.MgoSession.Clone()
		defer s.Close()
		db := s.DB(dummydata.Mongo.Database)
		ctx.Set(model.KeyDataDb, db)

		var post model.Post
		err := json.Unmarshal([]byte(dummydata.TestPost1), &post)
		So(err, ShouldBeNil)
		ctx.Set(model.KeyPost, post)

		PutPost(ctx)
		So(len(ctx.Errors), ShouldEqual, 0)

		posts := []model.Post{}
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
		So(err, ShouldBeNil)
		So(len(posts), ShouldEqual, 1)

		cleanMongoDB(t)

		posts2 := []model.Post{}
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
		So(err, ShouldBeNil)
		So(len(posts2), ShouldEqual, 0)

	})
}

// Convey("MultipleInserts", t, func() {
// 	s := dummydata.Session.Clone()
// 	defer s.Close()
// 	db := s.DB(dummydata.Mongo.Database)

// 	ctx, _ := mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	PutPost(ctx)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	ctx2, _ := mockContext(dummydata.TestPost2)
// 	ctx2.Set(model.KeyDb, db)
// 	PutPost(ctx2)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	posts := []model.Post{}
// 	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 	So(err, ShouldBeNil)
// 	So(len(posts), ShouldEqual, 2)

// 	cleanMongoDB(t)

// 	posts2 := []model.Post{}
// 	err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
// 	So(err, ShouldBeNil)
// 	So(len(posts2), ShouldEqual, 0)
// })

// Convey("Update", t, func() {
// 	s := dummydata.Session.Clone()
// 	defer s.Close()
// 	db := s.DB(dummydata.Mongo.Database)

// 	ctx, _ := mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	PutPost(ctx)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	// Update description
// 	{
// 		posts := []model.Post{}
// 		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 		So(err, ShouldBeNil)
// 		So(len(posts), ShouldEqual, 1)

// 		updatePost := posts[0]
// 		newDesc := "A slightly different description"
// 		updatePost.Desc = newDesc
// 		b, err := json.Marshal(updatePost)
// 		So(err, ShouldBeNil)

// 		ctx, _ := mockContext(string(b))
// 		ctx.Set(model.KeyDb, db)
// 		PutPost(ctx)
// 		So(len(ctx.Errors), ShouldEqual, 0)

// 		posts = []model.Post{}
// 		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 		So(err, ShouldBeNil)
// 		So(len(posts), ShouldEqual, 1)
// 		So(posts[0].Desc, ShouldEqual, newDesc)

// 	}

// 	// Try to update path => should fail
// 	{
// 		posts := []model.Post{}
// 		err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 		So(err, ShouldBeNil)
// 		So(len(posts), ShouldEqual, 1)

// 		updatePost := posts[0]
// 		newSlug := "a-slightly-different-slug"
// 		updatePost.Path = newSlug
// 		b, err := json.Marshal(updatePost)
// 		So(err, ShouldBeNil)

// 		ctx, _ := mockContext(string(b))
// 		ctx.Set(model.KeyDb, db)
// 		PutPost(ctx)
// 		So(len(ctx.Errors), ShouldEqual, 1)

// 		posts = []model.Post{}
// 		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 		So(err, ShouldBeNil)
// 		So(len(posts), ShouldEqual, 1)
// 		So(posts[0].Path, ShouldEqual, "simple-test")
// 	}

// 	cleanMongoDB(t)
// })

// Convey("TestPath", t, func() {

// 	s := dummydata.Session.Clone()
// 	defer s.Close()
// 	db := s.DB(dummydata.Mongo.Database)

// 	ctx, _ := mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	PutPost(ctx)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	// try to add same again
// 	ctx, _ = mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	PutPost(ctx)
// 	So(len(ctx.Errors), ShouldEqual, 1)

// 	posts := []model.Post{}
// 	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 	So(err, ShouldBeNil)
// 	So(len(posts), ShouldEqual, 1)

// 	exist := doesPathExist("simple-test", db)
// 	So(exist, ShouldBeTrue)

// 	exist = doesPathExist("simple-test-67", db)
// 	So(exist, ShouldBeFalse)

// 	cleanMongoDB(t)
// })

// Convey("Queries", t, func() {
// 	s := dummydata.Session.Clone()
// 	defer s.Close()
// 	db := s.DB(dummydata.Mongo.Database)

// 	ctx, _ := mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	PutPost(ctx)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	ctx2, _ := mockContext(dummydata.TestPost2)
// 	ctx2.Set(model.KeyDb, db)
// 	PutPost(ctx2)
// 	So(len(ctx.Errors), ShouldEqual, 0)

// 	posts := []model.Post{}
// 	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 	So(err, ShouldBeNil)
// 	So(len(posts), ShouldEqual, 2)

// 	cleanMongoDB(t)

// 	posts2 := []model.Post{}
// 	err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
// 	So(err, ShouldBeNil)
// 	So(len(posts2), ShouldEqual, 0)
// })

// Convey("TestDates", t, func() {

// 	s := dummydata.Session.Clone()
// 	defer s.Close()
// 	db := s.DB(dummydata.Mongo.Database)

// 	ctx, _ := mockContext(dummydata.TestPost1)
// 	ctx.Set(model.KeyDb, db)
// 	t1 := time.Now().Unix()
// 	time.Sleep(1 * time.Second)
// 	PutPost(ctx)
// 	time.Sleep(1 * time.Second)
// 	So(len(ctx.Errors), ShouldEqual, 0)
// 	t2 := time.Now().Unix()

// 	posts := []model.Post{}
// 	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
// 	So(err, ShouldBeNil)
// 	So(len(posts), ShouldEqual, 1)

// 	fp := posts[0]

// 	So(fp.Date, ShouldBeGreaterThan, t1)
// 	So(fp.Date, ShouldBeLessThan, t2)

// 	cleanMongoDB(t)
// })

// }
