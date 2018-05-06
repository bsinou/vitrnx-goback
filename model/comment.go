package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	// CommentCollection holds the name of the comment collection.
	CommentCollection = "comments"
)

// Comment model.
type Comment struct {
	ID       bson.ObjectId `json:"id,omitempty"  bson:"id,omitempty"`
	Date     int64         `json:"date,omitempty"`
	ParentID string        `json:"parentId" binding:"required" bson:"parentId"`
	AuthorID string        `json:"authorId" bson:"authorId"`
	Author   string        `json:"author" bson:"author"`
	Body     string        `json:"body" bson:"body"`

	CreatedOn int64  `json:"createdOn,omitempty" bson:"createdOn"`
	UpdatedOn int64  `json:"updatedOn,omitempty" bson:"updatedOn"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
