package model

import "gopkg.in/mgo.v2/bson"

const (
	// PostCollection holds the name of the post collection
	PostCollection = "posts"
)

// Post model
type Post struct {
	ID     bson.ObjectId `json:"id,omitempty" bson:"id,omitempty"`
	Path   string        `json:"path" binding:"required" bson:"path"`
	Title  string        `json:"title" binding:"required" bson:"title"`
	Author string        `json:"author" bson:"author"`
	Tags   string        `json:"tags" bson:"tags"`
	Desc   string        `json:"desc" bson:"desc"`
	Body   string        `json:"body" bson:"body"`

	CreatedOn int64 `json:"createdOn" bson:"createdOn"`
	UpdatedOn int64 `json:"updatedOn" bson:"updatedOn"`
}
