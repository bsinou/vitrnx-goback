package mongodb

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/bsinou/vitrnx-goback/conf"
)

var (
	// Session stores mongo session
	Session *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo
)

const (
	// defaultMongoURL is the failover URL for new empty instances
	defaultMongoURL = "mongodb://localhost:27017/vitrnx_store"
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
// 	Session = sy
// 	Mongo = mongoy
// }

// InitMongoConnection initialises a session with R/W privileges.
func InitMongoConnection() {
	uri := os.Getenv("MONGODB_URL")

	if len(uri) == 0 {
		uri = "mongodb://localhost:27017/" + conf.VitrnxInstanceID + "_store"
	}

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
	fmt.Println("Connected to", uri)
	Session = s
	Mongo = mongo
}
