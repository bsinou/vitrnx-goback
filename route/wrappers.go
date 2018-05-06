package route

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/auth"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"
	"github.com/bsinou/vitrnx-goback/mongodb"
)

// loggingHandler simply logs every request to stdout
func loggingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Got a  %s request at %s \n", c.Request.Method, c.Request.URL.String())

		t1 := time.Now()
		c.Next()
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), t2.Sub(t1))
	}
}

// checkCredentials calls the authentication API to verify the token
func checkCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Real check
		jwt := c.Request.Header.Get(model.KeyAuth)
		if jwt == "" {
			// No JWT cannot continue
			c.JSON(401, gin.H{"error": "No JWT token provided, please log out and back again"})
			c.Abort()
			return
		}

		err := auth.CheckCredentialAgainstFireBase(c, jwt)
		if err != nil {
			c.JSON(401, gin.H{"error": "Could not validate token: " + err.Error()})
			c.Abort()
			return
		}

		// Useless
		// c.Next()
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Setting CORS policies\n")

		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			// Headers have been set, no need to go further
			c.Abort()
			return
		}

		// 	NOTE: this is OK:
		// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 	This is *NOT* OK:  the second 'Authorization' line erase the first and Content-Type is not an authorized header anymore
		// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
	}
}

// checkCredentials calls the authentication API to verify the token
func checkCredentialsForUserCreation() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Checking credentials B4 user creation\n")

		// Real check
		jwt := c.Request.Header.Get(model.KeyAuth)

		err := auth.CheckCredentialAgainstFireBase(c, jwt)
		if err != nil {
			log.Printf("firebase credentials validation failed\n")
			c.JSON(503, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID := c.MustGet(model.KeyUserID).(string)
		var user model.User
		c.Bind(&user)
		// TODO enhance this will only be true upon user creation
		if userID != user.UserID {
			log.Printf("user ID differ: %s vs %s \n", userID, user.UserID)
			c.JSON(503, gin.H{"error": "Unauthorized"})
			c.Abort()
		}

		c.Set(model.KeyUser, user)

	}
}

// unmarshallPost retrieves the post from the context and store it again as a model.Post, for POST and MOVE requests.
// For delete request we rather use the passed path
// This must be done because the post can be retrieve only once from the request payload
func unmarshallPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		rm := c.Request.Method
		if rm == "POST" || rm == "MOVE" {
			post := model.Post{}
			err := c.Bind(&post)
			if err != nil {
				fmt.Printf("Could not bind post %v\n", err)
				c.Error(err)
				c.Abort()
				return
			}

			if post.Path == "" {
				err = fmt.Errorf("path is required, could not proceed")
				fmt.Println(err.Error())
				c.Error(err)
				c.Abort()
				return
			}
			c.Set(model.KeyPost, post)
		} else if rm == "DELETE" {
			post := model.Post{}
			path := c.Param(model.KeyPath)
			pathQuery := bson.M{model.KeyPath: path}
			db := c.MustGet(model.KeyDb).(*mgo.Database)
			err := db.C(model.PostCollection).Find(pathQuery).One(&post)
			if err != nil {
				fmt.Printf("Could not find post to delete with path %s: %s\n", path, err.Error())
				c.Error(err)
				c.Abort()
				return
			}
			c.Set(model.KeyPost, post)
		}
	}
}

// applyPostPolicies limit possible actions depending on user role
func applyPostPolicies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// shortcut
		rm := c.Request.Method

		// Get requests are filtered A POSTERIORI
		// if c.Request.Method == "GET" && !contains(claims, model.PolicyCanRead) {
		// 	c.JSON(503, gin.H{"error": "Unauthorized"})
		// 	c.Abort()
		// } else

		if rm == "POST" || rm == "DELETE" {
			userID := c.MustGet(model.KeyUserID).(string)
			roles := c.MustGet(model.KeyRoles).([]string)
			post := c.MustGet(model.KeyPost).(model.Post)

			switch rm {
			case "POST":
				// creation
				if post.ID.Hex() == "" && !(contains(roles, model.RoleModerator) || contains(roles, model.RoleEditor)) {
					c.JSON(403, gin.H{"error": "you don't have sufficient permission to create a post"})
					c.Abort()
					return
				}
				// update only by the post author or a moderator
				if !(contains(roles, model.RoleModerator) || post.AuthorID == userID) {
					c.JSON(403, gin.H{"error": "you don't have sufficient permission to update this post"})
					c.Abort()
					return
				}

				break
			case "DELETE":
				// delete only by the post author or a moderator
				if !(contains(roles, model.RoleModerator) || post.AuthorID == userID) {
					c.JSON(403, gin.H{"error": "you don't have sufficient permission to delete this post"})
					c.Abort()
					return
				}
			}
		}
	}
}

// Connect middleware clones the database session for each request and
// makes the `db` object available for each handler
func Connect(storeType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Connecting store %s \n", storeType)

		if storeType == model.StoreTypeMgo {
			s := mongodb.Session.Clone()
			defer s.Close()
			c.Set(model.KeyDb, s.DB(mongodb.Mongo.Database))
		} else if storeType == model.StoreTypeGorm {
			db := gorm.GetConnection()
			defer db.Close()
			c.Set(model.KeyDb, db)
		}
		// Next must be explicitely called here
		// so that the db session is released *AFTER* next handlers processing
		c.Next()
	}
}

// // ErrorHandler is a middleware to handle errors encountered during requests
// func ErrorHandler(c *gin.Context) {
// 	c.Next()

// 	// TODO: Handle it in a better way
// 	if len(c.Errors) > 0 {
// 		c.HTML(http.StatusBadRequest, "400", gin.H{
// 			"errors": c.Errors,
// 		})
// 	}
// }

func contains(arr []string, val string) bool {
	for _, currVal := range arr {
		if currVal == val {
			return true
		}
	}
	return false
}
