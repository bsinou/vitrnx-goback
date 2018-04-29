package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
)

// ListPosts retrieves all posts, potentially filtered by passed tag value.
func ListPosts(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	queryTag := c.Query("tag")
	posts := []model.Post{}

	var err error
	if queryTag == "" {
		err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	} else {
		query := bson.M{model.KeyTags: queryTag}
		err = db.C(model.PostCollection).Find(query).Sort("-updatedOn").All(&posts)
	}
	if err != nil {
		c.Error(err)
	}

	// fmt.Printf("Retrieved %d posts\n", len(posts))
	// if len(posts) > 0 {
	// 	fmt.Printf("Id of first retrieved posts: %v \n", posts[0])
	// }

	c.JSON(200, posts)
}

/* CRUD */

// PutPost simply creates or updates a post in the document repository.
func PutPost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	fmt.Printf("In put post\n")

	post := model.Post{}
	err := c.Bind(&post)
	if err != nil {
		fmt.Printf("Could not bind post %v\n", err)
		c.Error(err)
		return
	}

	if post.Path == "" {
		err = fmt.Errorf("path is required, could not upsert")
		fmt.Println(err.Error())
		c.Error(err)
		return
	}

	posts := db.C(model.PostCollection)
	creation := false

	if post.ID.Hex() == "" {
		// Check path unicity
		if doesPathExist(post.Path, db) {
			err = fmt.Errorf("could not create: a post already exist at %s", post.Path)
			c.Error(err)
			return
		}

		creation = true
		// Set creation info
		post.ID = bson.NewObjectId()
		post.Date = time.Now()
		post.Author = c.MustGet(model.KeyUserName).(string)
		post.CreatedOn = time.Now().Unix()
	} else {
		// Prevent move
		var oldPost model.Post
		query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}
		err := db.C(model.PostCollection).Find(query).One(&oldPost)
		if err != nil {
			fmt.Printf("update failed: could not find post with id %s, %v\n", post.ID, err)
			c.Error(err)
			return
		}

		if oldPost.Path != post.Path {
			fmt.Printf("different paths  %s != %s \n", oldPost.Path, post.Path)
			c.Error(fmt.Errorf("update failed: cannot modify path for %s", oldPost.Path))
			return
		}
	}

	// Always update the update (...) info
	post.UpdatedOn = time.Now().Unix()
	post.UpdatedBy = c.MustGet(model.KeyUserName).(string)

	if creation {
		err = posts.Insert(post)
		if err != nil {
			fmt.Printf("Insert failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(201, gin.H{"post": post})
	} else {
		query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}
		err = posts.Update(query, post)
		if err != nil {
			fmt.Printf("Update failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(200, gin.H{"post": post})
	}
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

	c.JSON(201, gin.H{"post": post, "claims": `{"canManage": "true"}`})
}

// DeletePost definitively removes a post from the repository
func DeletePost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	post := model.Post{}
	err := c.Bind(&post)
	if err != nil {
		fmt.Printf("Could not bind post %v\n", err)
		c.Error(err)
		return
	}
	path := post.Path
	query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}

	err = db.C(model.PostCollection).Remove(query)
	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("post at %s has been deleted", path)})

}

/* Helper functions */

func doesPathExist(path string, db *mgo.Database) bool {
	post := model.Post{}
	pathQuery := bson.M{model.KeyPath: path}
	// Maybe add a check to insure unicity?
	err := db.C(model.PostCollection).Find(pathQuery).One(&post)
	return err == nil
}

// func updatePost(c *gin.Context) {
// 	db := c.MustGet(model.KeyDb).(*mgo.Database)

// 	updatedPost := model.Post{}
// 	err := c.Bind(&updatedPost)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	query := bson.M{"id": bson.ObjectIdHex(c.Param("id"))}
// 	doc := bson.M{
// 		"path":      updatedPost.Path,
// 		"title":     updatedPost.Title,
// 		"body":      updatedPost.Body,
// 		"updatedOn": time.Now().UnixNano() / int64(time.Millisecond),
// 	}
// 	err = db.C(model.PostCollection).Update(query, doc)
// 	if err != nil {
// 		c.Error(err)
// 	}
// 	c.JSON(201, gin.H{"success": updatedPost})

// 	c.Redirect(http.StatusMovedPermanently, "/posts")
// }
