package gorm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	. "github.com/smartystreets/goconvey/convey"
)

type Blog struct {
	ID      uint `gorm:"primary_key"`
	Subject string
	Body    string
	Tags    []Tag `gorm:"many2many:blog_tags;"`
}

type Tag struct {
	ID    string `gorm:"primary_key"`
	Value string
	// Blogs []*Blog `gorm:"many2many:blog_tags"`
}

func compareTags(tags []Tag, contents []string) bool {
	var tagContents []string
	for _, tag := range tags {
		tagContents = append(tagContents, tag.Value)
	}
	sort.Strings(tagContents)
	sort.Strings(contents)
	return reflect.DeepEqual(tagContents, contents)
}

func createDb() (*gorm.DB, string) {
	fname := fmt.Sprintf("/sqlite-%d.db", time.Now().UnixNano())
	path := os.TempDir() + fname

	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot open sqlite at %s, : %s\n", dbFileAbsPath, err))
	}
	// db.LogMode(true) // Display SQL queries

	db.DropTable(&Blog{}, &Tag{})
	db.DropTable("blog_tags")
	db.CreateTable(&Blog{}, &Tag{})

	return db, path
}

func cleanDb(db *gorm.DB, path string) {
	db.Close()
	os.Remove(path)
}

func TestManyToManyWithMultiPrimaryKeys(t *testing.T) {

	Convey("Test many to many associations: ", t, func() {
		db, path := createDb()
		defer cleanDb(db, path)

		blog := Blog{
			Subject: "subject",
			Body:    "body",
			Tags: []Tag{
				{ID: "TAG_1", Value: "tag1"},
				{ID: "TAG_2", Value: "tag2"},
			},
		}
		db.Save(&blog)
		if !compareTags(blog.Tags, []string{"tag1", "tag2"}) {
			t.Errorf("Blog should has two tags")
		}

		blog2 := Blog{
			Subject: "subject2",
			Body:    "body2",
			Tags: []Tag{
				{ID: "TAG_1", Value: "tag1"},
				{ID: "TAG_2", Value: "tag2"},
			},
		}
		db.Save(&blog2)

		var tags []Tag
		db.Find(&tags)
		So(len(tags), ShouldEqual, 2)

		var blogs []Blog
		db.Preload("Tags").Find(&blogs)
		for _, currBlog := range blogs {
			tstStr, err := json.Marshal(&currBlog)
			So(err, ShouldBeNil)
			fmt.Println(string(tstStr) + "\n")
		}
	})

	Convey("Test tag value update: ", t, func() {
		db, path := createDb()
		defer cleanDb(db, path)

		blog := Blog{
			Subject: "subject",
			Body:    "body",
			Tags: []Tag{
				{ID: "TAG_1", Value: "tag2"},
				{ID: "TAG_2", Value: "tag2"},
			},
		}
		db.Save(&blog)
		if !compareTags(blog.Tags, []string{"tag1", "tag2"}) {
			t.Errorf("Blog should have two tags")
		}

		blog2 := Blog{
			Subject: "subject2",
			Body:    "body2",
			Tags: []Tag{
				{ID: "TAG_1"},
				{ID: "TAG_2"},
			},
		}
		db.Set("gorm:association_autoupdate", false).Save(&blog2)

		var tags []Tag
		db.Find(&tags)
		So(len(tags), ShouldEqual, 2)

		var blogs []Blog
		db.Preload("Tags").Find(&blogs)
		for _, currBlog := range blogs {
			tstStr, err := json.Marshal(&currBlog)
			So(err, ShouldBeNil)
			fmt.Println(string(tstStr) + "\n")
		}
	})

	// // Append
	// var tag3 = &Tag{Locale: "ZH", Value: "tag3"}
	// db.Model(&blog).Association("Tags").Append([]*Tag{tag3})
	// if !compareTags(blog.Tags, []string{"tag1", "tag2", "tag3"}) {
	// 	t.Errorf("Blog should has three tags after Append")
	// }

	// if db.Model(&blog).Association("Tags").Count() != 3 {
	// 	t.Errorf("Blog should has three tags after Append")
	// }

	// var tags []Tag
	// db.Model(&blog).Related(&tags, "Tags")
	// if !compareTags(tags, []string{"tag1", "tag2", "tag3"}) {
	// 	t.Errorf("Should find 3 tags with Related")
	// }

	// var blog1 Blog
	// db.Preload("Tags").Find(&blog1)
	// if !compareTags(blog1.Tags, []string{"tag1", "tag2", "tag3"}) {
	// 	t.Errorf("Preload many2many relations")
	// }

	// // Replace
	// var tag5 = &Tag{Locale: "ZH", Value: "tag5"}
	// var tag6 = &Tag{Locale: "ZH", Value: "tag6"}
	// db.Model(&blog).Association("Tags").Replace(tag5, tag6)
	// var tags2 []Tag
	// db.Model(&blog).Related(&tags2, "Tags")
	// if !compareTags(tags2, []string{"tag5", "tag6"}) {
	// 	t.Errorf("Should find 2 tags after Replace")
	// }

	// if db.Model(&blog).Association("Tags").Count() != 2 {
	// 	t.Errorf("Blog should has three tags after Replace")
	// }

	// // Delete
	// db.Model(&blog).Association("Tags").Delete(tag5)
	// var tags3 []Tag
	// db.Model(&blog).Related(&tags3, "Tags")
	// if !compareTags(tags3, []string{"tag6"}) {
	// 	t.Errorf("Should find 1 tags after Delete")
	// }

	// if db.Model(&blog).Association("Tags").Count() != 1 {
	// 	t.Errorf("Blog should has three tags after Delete")
	// }

	// db.Model(&blog).Association("Tags").Delete(tag3)
	// var tags4 []Tag
	// db.Model(&blog).Related(&tags4, "Tags")
	// if !compareTags(tags4, []string{"tag6"}) {
	// 	t.Errorf("Tag should not be deleted when Delete with a unrelated tag")
	// }

	// // Clear
	// db.Model(&blog).Association("Tags").Clear()
	// if db.Model(&blog).Association("Tags").Count() != 0 {
	// 	t.Errorf("All tags should be cleared")
	// }
	// })
}
