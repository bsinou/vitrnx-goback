package model

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	PresenceCollection = "presence"
)

// Presence model.
type Presence struct {
	ID     bson.ObjectId `json:"id,omitempty"  bson:"id,omitempty"`
	UserID string        `json:"userId" binding:"required" bson:"userId"`
	HasCar bool          `json:"hasCar" bson:"hasCar"`
	D1     int8          `json:"d1" bson:"d1"`
	D1C    int8          `json:"d1c" bson:"d1c"`
	D2     int8          `json:"d2" bson:"d2"`
	D2C    int8          `json:"d2c" bson:"d2c"`
	D3     int8          `json:"d3" bson:"d3"`
	D3C    int8          `json:"d3c" bson:"d3c"`
	D4     int8          `json:"d4" bson:"d4"`
	D4C    int8          `json:"d4c" bson:"d4c"`
	D5     int8          `json:"d5" bson:"d5"`
	D5C    int8          `json:"d5c" bson:"d5c"`
	D6     int8          `json:"d6" bson:"d6"`
	D6C    int8          `json:"d6c" bson:"d6c"`
	D7     int8          `json:"d7" bson:"d7"`
	D7C    int8          `json:"d7c" bson:"d7c"`
	D8     int8          `json:"d8" bson:"d8"`
	D8C    int8          `json:"d8c" bson:"d8c"`

	UpdatedOn int64  `json:"updatedOn,omitempty" bson:"updatedOn"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
