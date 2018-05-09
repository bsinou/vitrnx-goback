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

func addUserMeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.WithUserMeta(c)
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

// TODO mettre en lower case....

// Connect middleware makes the various `db` objects available for each handler
func Connect() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Connecting stores\n")

		// User management
		db := gorm.GetConnection()
		defer db.Close()
		c.Set(model.KeyUserDb, db)

		// Data
		s := mongodb.Session.Clone()
		defer s.Close()
		c.Set(model.KeyDataDb, s.DB(mongodb.Mongo.Database))

		// Next must be explicitely called here
		// so that the db session is released *AFTER* next handlers processing
		c.Next()
	}
}

/* USER */

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

/* POSTS */

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
			db := c.MustGet(model.KeyDataDb).(*mgo.Database)
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
			roles := c.MustGet(model.KeyUserRoles).([]string)
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

/* COMMENTS */

// unmarshallComment retrieves the comment from the context and store it again as a model.Comment, for POST and MOVE requests.
// For delete request we use the passed id
func unmarshallComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		rm := c.Request.Method
		if rm == "POST" || rm == "MOVE" {
			comment := model.Comment{}
			err := c.Bind(&comment)
			if err != nil {
				fmt.Printf("Could not bind comment %v\n", err)
				c.Error(err)
				c.Abort()
				return
			}

			if comment.ParentID == "" {
				err = fmt.Errorf("parentId is required, could not upsert")
				fmt.Println(err.Error())
				c.Error(err)
				return
			}
			c.Set(model.KeyComment, comment)
		} else if rm == "DELETE" {
			db := c.MustGet(model.KeyDataDb).(*mgo.Database)
			comment := model.Comment{}
			id := c.Param(model.KeyMgoID)
			idQuery := bson.M{model.KeyMgoID: bson.ObjectIdHex(id)}
			err := db.C(model.CommentCollection).Find(idQuery).One(&comment)
			if err != nil {
				fmt.Printf("Could not find comment to delete with id %s: %s\n", id, err.Error())
				c.Error(err)
				c.Abort()
				return
			}
			c.Set(model.KeyComment, comment)
		}
	}
}

// applyCommentPolicies limit possible actions depending on user role
func applyCommentPolicies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// shortcut
		rm := c.Request.Method

		// Note: GET requests are filtered A POSTERIORI
		if rm == "POST" || rm == "DELETE" {
			userID := c.MustGet(model.KeyUserID).(string)
			roles := c.MustGet(model.KeyUserRoles).([]string)
			comment := c.MustGet(model.KeyComment).(model.Comment)

			switch rm {
			case "POST":
				// creation
				// TODO prevent addition when guest
				// if comment.ID.Hex() == "" {
				// 	// && !(contains(roles, model.RoleModerator) || contains(roles, model.RoleEditor)) {
				// 	c.JSON(403, gin.H{"error": "you don't have sufficient permission to add a comment "})
				// 	c.Abort()
				// 	return
				// }
				// update only by the post author or a moderator
				if !(contains(roles, model.RoleModerator) || comment.AuthorID == userID) {
					c.JSON(403, gin.H{"error": "you don't have sufficient permission to update this comment"})
					c.Abort()
					return
				}

				break
			case "DELETE":
				// only the comment author or a moderator can delete
				if !(contains(roles, model.RoleModerator) || comment.AuthorID == userID) {
					c.JSON(403, gin.H{"error": "you don't have sufficient permission to delete this post"})
					c.Abort()
					return
				}
			}
		}
	}
}

/* HELPER FUNCTIONS */

func contains(arr []string, val string) bool {
	for _, currVal := range arr {
		if currVal == val {
			return true
		}
	}
	return false
}
