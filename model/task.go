package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	TaskCollection = "tasks"
)

// Task model.
type Task struct {
	ID         bson.ObjectId `json:"id,omitempty"  bson:"id,omitempty"`
	Flags      int           `json:"flags"  bson:"flags"`
	CategoryId string        `json:"categoryId"  bson:"categoryId"`
	ManagerID  string        `json:"managerId" bson:"managerId"`
	Manager    string        `json:"manager" bson:"manager"`
	Desc       string        `json:"desc" bson:"desc"`

	CreationDate int64  `json:"creationDate" bson:"creationDate"`
	CloseDate    int64  `json:"closeDate" bson:"closeDate"`
	CancelDate   int64  `json:"cancelDate" bson:"cancelDate"`
	DueDate      int64  `json:"dueDate" bson:"dueDate"`
	UpdateDate   int64  `json:"updatedOn,omitempty" bson:"updatedOn"`
	UpdateBy     string `json:"updatedBy" bson:"updatedBy"`
}
