package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
	// "time"

	// "gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"

	. "github.com/smartystreets/goconvey/convey"
)

// FIXME check this
//  https://github.com/jinzhu/gorm/blob/021d7b33143de37b743d1cf660974e9c8d3f80ea/multi_primary_keys_test.go

func TestUserDao_CRUD(t *testing.T) {

	fname := fmt.Sprintf("/sqlite-%d.db", time.Now().Unix())
	path := os.TempDir() + fname
	gorm.InitGormTestRepo(path)
	defer func() {
		os.Remove(path)
	}()

	Convey("Test Roles mapping: ", t, func() {

		user1 := getPreparedUser("john", "admin")
		printUser(user1)

		user2 := getPreparedUser("jane", "admin")
		printUser(user2)

		db := gorm.GetConnection()
		defer db.Close()

		db.Save(&user1)
		db.Model(&user1).Association("Roles").Append(model.Role{RoleID: "admin"})

		db.Create(&user2)
		// db.Save(&user2)

		var roles []model.Role
		db.Find(&roles)
		So(len(roles), ShouldEqual, 4)

		var users []model.User
		db.Preload("Roles").Find(&users)
		So(len(users), ShouldEqual, 2)

		for _, user := range users {
			tstStr, err := json.Marshal(&user)
			So(err, ShouldBeNil)
			fmt.Println(string(tstStr) + "\n")
		}
	})

	Convey("Test model: ", t, func() {
		role0 := model.Role{
			RoleID: "EDITOR",
			Label:  "Editor",
		}
		role1 := model.Role{
			RoleID: "ADMIN",
			Label:  "Admin",
		}

		user0 := model.User{
			UserID: "test",
			Roles:  []model.Role{role0, role1},
			Name:   "John",
			Email:  "john@example.com",
		}

		tstStr, err := json.Marshal(&user0)
		So(err, ShouldBeNil)
		fmt.Println(string(tstStr))
	})

	// Convey("Test Roles: ", t, func() {

	// 	ctx, _ := mockContext()
	// 	db := gorm.GetConnection()
	// 	defer db.Close()
	// 	ctx.Set(model.KeyDb, db)

	// 	var user model.User
	// 	err := json.Unmarshal([]byte(dummydata.TestUser1), &user)
	// 	So(err, ShouldBeNil)
	// 	ctx.Set(model.KeyUser, user)

	// 	PutUser(ctx)
	// 	So(len(ctx.Errors), ShouldEqual, 0)

	// 	ctx, _ = mockContext()
	// 	ctx.Set(model.KeyDb, db)

	// 	var user2 model.User
	// 	err = json.Unmarshal([]byte(dummydata.TestUser2), &user2)
	// 	So(err, ShouldBeNil)
	// 	ctx.Set(model.KeyUser, user2)

	// 	PutUser(ctx)
	// 	So(len(ctx.Errors), ShouldEqual, 0)

	// 	var users []model.User
	// 	db.Find(&users)
	// 	So(len(users), ShouldEqual, 2)

	// 	var roles []model.Role
	// 	db.Find(&roles)
	// 	So(len(roles), ShouldEqual, 2)

	// 	for _, user := range users {
	// 		tstStr, err := json.Marshal(&user)
	// 		So(err, ShouldBeNil)
	// 		fmt.Println(string(tstStr))
	// 	}

	// 	// posts := []model.Post{}
	// 	// err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	// 	// So(err, ShouldBeNil)
	// 	// So(len(posts), ShouldEqual, 1)

	// 	// cleanMongoDB(t)

	// 	// posts2 := []model.Post{}
	// 	// err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts2)
	// 	// So(err, ShouldBeNil)
	// 	// So(len(posts2), ShouldEqual, 0)

	// })

}

func printUser(user *model.User) {
	tstStr, err := json.Marshal(user)
	So(err, ShouldBeNil)
	fmt.Println(string(tstStr))
}
