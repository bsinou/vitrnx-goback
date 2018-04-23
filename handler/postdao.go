package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
)

// New post
func NewPost(c *gin.Context) {
	post := model.Post{}

	c.HTML(http.StatusOK, "posts/form", gin.H{
		"title": "New post",
		"post":  post,
	})
}

// Create a post
func CreatePost(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	post := model.Post{}
	err := c.Bind(&post)
	if err != nil {
		c.Error(err)
		return
	}

	err = db.C(model.PostCollection).Insert(post)
	if err != nil {
		// c.JSON(422, gin.H{"error": "Not enough info"})
		c.Error(err)
	}

	c.JSON(201, gin.H{"success": post})
}

// Edit a post
func EditPost(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	post := model.Post{}
	oID := bson.ObjectIdHex(c.Param("id"))
	err := db.C(model.PostCollection).FindId(oID).One(&post)
	if err != nil {
		c.Error(err)
	}

	c.JSON(201, gin.H{"success": post})

	// c.HTML(http.StatusOK, "posts/form", gin.H{
	// 	"title": "Edit post",
	// 	"post":  post,
	// })
}

// List all posts
func ListPosts(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	fmt.Println("Listing posts")

	posts := []model.Post{}
	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	if err != nil {
		c.Error(err)
	}
	fmt.Printf("Retrieved %d posts\n", len(posts))

	c.JSON(200, posts)
}

// Update an post
func UpdatePost(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)

	post := model.Post{}
	err := c.Bind(&post)
	if err != nil {
		c.Error(err)
		return
	}

	query := bson.M{"id": bson.ObjectIdHex(c.Param("id"))}
	doc := bson.M{
		"path":      post.Path,
		"title":     post.Title,
		"body":      post.Body,
		"updatedOn": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = db.C(model.PostCollection).Update(query, doc)
	if err != nil {
		c.Error(err)
	}
	c.JSON(201, gin.H{"success": post})

	c.Redirect(http.StatusMovedPermanently, "/posts")
}

// Delete a post
func DeletePost(c *gin.Context) {
	db := c.MustGet("db").(*mgo.Database)
	query := bson.M{"id": bson.ObjectIdHex(c.Param("id"))}
	err := db.C(model.PostCollection).Remove(query)
	if err != nil {
		c.Error(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/posts")
}
