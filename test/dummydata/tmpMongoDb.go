package dummydata

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var (
	// MgoSession stores mongo session
	MgoSession *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo
)

const (
	// defaultMongoURL defines a failover connection forthe tests.
	defaultMongoURL = "mongodb://localhost:27017/test_store"
)

// init connects to mongodb
func init() {
	uri := os.Getenv("MONGODB_URL")

	if len(uri) == 0 {
		uri = defaultMongoURL
	}

	mongo, err := mgo.ParseURL(uri)
	s, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}
	s.SetSafe(&mgo.Safe{})
	fmt.Println("In dummydata, connected to ", uri)
	MgoSession = s
	Mongo = mongo
}
