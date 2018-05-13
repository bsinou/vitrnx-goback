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

// ListPosts retrieves all posts, potentially filtered by passed tag value.
func ListPosts(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	posts := []model.Post{}
	var err error

	// err = db.C(model.PostCollection).Find(nil).Sort("-updatedOn").All(&posts)
	// fmt.Println("Tag: " + queryTag)

	queryTag := c.Query("tag")
	if queryTag == "" {
		err = db.C(model.PostCollection).Find(nil).Sort("-createdOn").All(&posts)
	} else {
		query := bson.M{model.KeyTags: bson.RegEx{queryTag, ""}}
		err = db.C(model.PostCollection).Find(query).Sort("-createdOn").All(&posts)
	}
	if err != nil {
		c.Error(err)
	}

	// might be enhanced using https://stackoverflow.com/questions/37562873/most-idiomatic-way-to-select-elements-from-an-array-in-golang
	var updatedPosts []model.Post
	for _, post := range posts {
		query := bson.M{model.KeyParentID: bson.RegEx{post.Path, ""}}
		post.CommentCount, err = db.C(model.CommentCollection).Find(query).Count()
		updatedPosts = append(updatedPosts, post)
		fmt.Printf("retrieving post at %s with %d comments \n", post.Path, post.CommentCount)
		// TODO also filter by rights
	}

	c.JSON(200, gin.H{"posts": updatedPosts})
}

/* CRUD */

// PutPost simply creates or updates a post in the document repository.
func PutPost(c *gin.Context) {
	post := c.MustGet(model.KeyPost).(model.Post)
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	posts := db.C(model.PostCollection)
	creation := false

	if post.ID.Hex() == "" {
		// Check path unicity
		if doesPathExist(post.Path, db) {
			err := fmt.Errorf("could not create: a post already exist at %s", post.Path)
			c.Error(err)
			return
		}

		creation = true
		// Set creation info
		post.ID = bson.NewObjectId()
		post.Date = time.Now().Unix()
		post.AuthorID = c.MustGet(model.KeyUserID).(string)
		post.Author = c.MustGet(model.KeyUserDisplayName).(string)
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

		// // TODO clean this: date is lost on update
		// // fmt.Println("### Updating: received date: " + post.Date.Format("2006-01-02"))
		// // fmt.Println("### Updating: old date: " + oldPost.Date.Format("2006-01-02"))
		// post.Date = oldPost.Date
		// post.AuthorID = oldPost.AuthorID
		// post.Author = oldPost.Author
	}

	// Always update the update (...) info
	post.UpdatedOn = time.Now().Unix()
	post.UpdatedBy = c.MustGet(model.KeyUserID).(string)

	if creation {
		err := posts.Insert(post)
		if err != nil {
			fmt.Printf("Insert failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(201, gin.H{"post": post})
	} else {
		query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}
		err := posts.Update(query, post)
		if err != nil {
			fmt.Printf("Update failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(200, gin.H{"post": post})
	}
}

// ReadPost simply retrieves a post by path
func ReadPost(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	post := model.Post{}
	pathQuery := bson.M{model.KeyPath: c.Param(model.KeyPath)}
	err := db.C(model.PostCollection).Find(pathQuery).One(&post)
	if err != nil {
		c.Error(err)
		return
	}

	// TODO check if current user can see this post
	c.JSON(201, gin.H{"post": post})
}

// ListPostComments retrieves the comment for a given post
func ListPostComments(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	comments := []model.Comment{}

	path := c.Param(model.KeyPath)
	query := bson.M{model.KeyParentID: bson.RegEx{path, ""}}
	err := db.C(model.CommentCollection).Find(query).Sort("-date").All(&comments)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"comments": comments})
}

// DeletePost definitively removes a post from the repository
func DeletePost(c *gin.Context) {
	post := c.MustGet(model.KeyPost).(model.Post)
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	query := bson.M{"id": bson.ObjectIdHex(post.ID.Hex())}
	err := db.C(model.PostCollection).Remove(query)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"message": fmt.Sprintf("post at %s has been deleted", post.Path)})
	c.Redirect(http.StatusMovedPermanently, "/posts")
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
