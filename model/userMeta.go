package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	PresenceCollection = "presence"
)

// Presence model.
type Presence struct {
	ID        bson.ObjectId `json:"id,omitempty"  bson:"id,omitempty"`
	UserID    string        `json:"userId" binding:"required" bson:"userId"`
	IsComing  bool          `json:"isComing" bson:"isComing"`
	Transport string        `json:"transport" bson:"transport"`
	AdultNb   int8          `json:"adultNb" bson:"adultNb"`
	ChildNb   int8          `json:"childNb" bson:"childNb"`
	D1        bool          `json:"d1" bson:"d1"`
	D2        bool          `json:"d2" bson:"d2"`
	D3        bool          `json:"d3" bson:"d3"`
	D4        bool          `json:"d4" bson:"d4"`
	D5        bool          `json:"d5" bson:"d5"`
	D6        bool          `json:"d6" bson:"d6"`
	D7        bool          `json:"d7" bson:"d7"`
	D8        bool          `json:"d8" bson:"d8"`
	Comments  string        `json:"comments" bson:"comments"`

	UpdatedOn int64  `json:"updatedOn,omitempty" bson:"updatedOn"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
