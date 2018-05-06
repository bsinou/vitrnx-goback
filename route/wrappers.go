package route

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

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
		err := auth.CheckCredentialAgainstFireBase(c, jwt)
		if err != nil {
			// this is not enough, the list is still sent back.
			c.JSON(503, gin.H{"error": "Unauthorized"})
			// We have to explicitely abort the request
			c.Abort()
			return
		}

		fmt.Println("Authorized, about to forward...")
		// Useless
		// c.Next()
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
func applyPostPolicies() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet(model.KeyClaims).([]string)

		if c.Request.Method == "GET" && !contains(claims, model.PolicyCanRead) {
			c.JSON(503, gin.H{"error": "Unauthorized"})
			c.Abort()
		} else if c.Request.Method == "POST" && !contains(claims, model.PolicyCanEdit) {
			c.JSON(503, gin.H{"error": "Unauthorized"})
			c.Abort()
		} else if c.Request.Method == "DELETE" && !contains(claims, model.PolicyCanManage) {
			c.JSON(503, gin.H{"error": "Unauthorized"})
			c.Abort()
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
