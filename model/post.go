package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// PostCollection holds the name of the post collection
	PostCollection = "posts"
)

// Post model
type Post struct {
	// --> ID was not automatically generated and serialized in JSON when called 'id'
	// _id seems to be OK
	ID     bson.ObjectId `json:"id,omitempty"  bson:"_id,omitempty"`
	Date   time.Time     `json:"date,omitempty"`
	Path   string        `json:"path" binding:"required" bson:"path"`
	Title  string        `json:"title" binding:"required" bson:"title"`
	Author string        `json:"author" bson:"author"`
	Tags   string        `json:"tags" bson:"tags"`
	Desc   string        `json:"desc" bson:"desc"`
	Body   string        `json:"body" bson:"body"`

	CreatedOn int64  `json:"createdOn,omitempty" bson:"createdOn"`
	UpdatedOn int64  `json:"updatedOn,omitempty" bson:"updatedOn"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
