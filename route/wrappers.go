package route

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/gateway"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"
	"github.com/bsinou/vitrnx-goback/mongodb"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// loggingHandler simply logs every request to stdout
func loggingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Got a request\n")

		t1 := time.Now()
		c.Next()
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), t2.Sub(t1))
	}
}

// checkCredentials calls the authentication API to verify the token
func checkCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == "OPTIONS" {
			fmt.Println("Silently returning...")
			return
		}

		fmt.Println("Method: " + c.Request.Method)
		// for k, v := range c.Request.Header {
		// 	fmt.Printf("%s - %s\n", k, v)
		// }

		// Real check
		jwt := c.Request.Header.Get(model.KeyAuth)
		fmt.Println("Header key: " + jwt)
		err := gateway.CheckCredentialAgainstFireBase(c, jwt)
		if err != nil {
			// if c.Request.Header.Get(model.KeyAuth) != "AUTH TOKEN" {
			// TODO this is not enough, the list is still sent back.
			c.JSON(503, gin.H{"error": "Unauthorized"})
			return
		}

		fmt.Println("Authorized, about to forward...")
		c.Next()
	}
}

// Connect middleware clones the database session for each request and
// makes the `db` object available for each handler
func Connect(storeType string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if storeType == model.StoreTypeMgo {
			s := mongodb.Session.Clone()
			defer s.Close()
			c.Set(model.KeyDb, s.DB(mongodb.Mongo.Database))
		} else if storeType == model.StoreTypeGorm {
			db := gorm.GetConnection()
			defer db.Close()
			c.Set("db", db)
		}

		c.Next()

	}
}

// ErrorHandler is a middleware to handle errors encountered during requests
func ErrorHandler(c *gin.Context) {
	c.Next()

	// TODO: Handle it in a better way
	if len(c.Errors) > 0 {
		c.HTML(http.StatusBadRequest, "400", gin.H{
			"errors": c.Errors,
		})
	}
}
