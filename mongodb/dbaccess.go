package mongodb

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var (
	// Session stores mongo session
	Session *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	MongoDBUrl = "mongodb://localhost:27017/vitrnx_store"
)

// init connects to mongodb
// func init() {
// 	uri := os.Getenv("MONGODB_URL")

// 	if len(uri) == 0 {
// 		uri = MongoDBUrl
// 	}

// 	mongo, err := mgo.ParseURL(uri)
// 	s, err := mgo.Dial(uri)
// 	if err != nil {
// 		fmt.Printf("Can't connect to mongo, go error %v\n", err)
// 		panic(err.Error())
// 	}
// 	s.SetSafe(&mgo.Safe{})
// 	fmt.Println("Connected to", uri)
// 	Session = s
// 	Mongo = mongo
// }

func InitMongoConnection() {
	uri := os.Getenv("MONGODB_URL")

	if len(uri) == 0 {
		uri = MongoDBUrl
	}

	mongo, err := mgo.ParseURL(uri)
	s, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}
	s.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", uri)
	Session = s
	Mongo = mongo
}
