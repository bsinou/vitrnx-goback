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

// List all posts
func ListPosts(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	posts := []model.Post{}
	err := db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	if err != nil {
		c.Error(err)
	}
	fmt.Printf("Retrieved %d posts\n", len(posts))
	if len(posts) > 0 {
		fmt.Printf("Id of first retrieved posts: %v \n", posts[0])
	}

	c.JSON(200, posts)
}

/* CRUD */

// PutPost simply creates or updates a post in the document repository.
func PutPost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	post := model.Post{}
	err := c.Bind(&post)
	if err != nil {
		c.Error(err)
		return
	}

	fmt.Printf("Before updating the post, id: %s\n", post.ID)

	// Only on create for the time being
	post.Date = time.Now()
	post.Author = c.MustGet(model.KeyUserName).(string)
	post.CreatedOn = time.Now().Unix()

	// Overwrite with curr value
	// TODO add change check and versioning
	post.UpdatedOn = time.Now().Unix()
	post.UpdatedBy = c.MustGet(model.KeyUserName).(string)

	posts := db.C(model.PostCollection)

	info, err := posts.Upsert(nil, post)
	if err != nil {
		// c.JSON(422, gin.H{"error": "Not enough info"})
		c.Error(err)
	}

	if info.UpsertedId != nil {
		post.ID = info.UpsertedId.(bson.ObjectId)
	} else {
		fmt.Printf("No ID generated ... \n")
	}

	fmt.Printf("Post upserted with ID: %v\n", post.ID.String())
	fmt.Printf("Change info: %v\n", info)
	c.JSON(201, gin.H{"post": post})
}

// ReadPost simply retrieves a post by path
func ReadPost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	post := model.Post{}
	pathQuery := bson.M{model.KeyPath: c.Param(model.KeyPath)}
	err := db.C(model.PostCollection).Find(pathQuery).One(&post)
	if err != nil {
		c.Error(err)
	}

	c.JSON(201, gin.H{"post": post})
}

// DeletePost definitively removes a post from the repository
func DeletePost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)
	query := bson.M{"id": bson.ObjectIdHex(c.Param("id"))}
	err := db.C(model.PostCollection).Remove(query)
	if err != nil {
		c.Error(err)
	}
	c.Redirect(http.StatusMovedPermanently, "/posts")
}

func updatePost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	updatedPost := model.Post{}
	err := c.Bind(&updatedPost)
	if err != nil {
		c.Error(err)
		return
	}

	query := bson.M{"id": bson.ObjectIdHex(c.Param("id"))}
	doc := bson.M{
		"path":      updatedPost.Path,
		"title":     updatedPost.Title,
		"body":      updatedPost.Body,
		"updatedOn": time.Now().UnixNano() / int64(time.Millisecond),
	}
	err = db.C(model.PostCollection).Update(query, doc)
	if err != nil {
		c.Error(err)
	}
	c.JSON(201, gin.H{"success": updatedPost})

	c.Redirect(http.StatusMovedPermanently, "/posts")
}
