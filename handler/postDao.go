package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	jgorm "github.com/jinzhu/gorm"

	// Use SQLite
	_ "github.com/mattn/go-sqlite3"

	"github.com/bsinou/vitrnx-goback/model"
)

// For the record defined model

// Date      int64  `json:"date,omitempty"`
// Slug      string `json:"slug"`
// Title     string `json:"title"`
// AuthorID  string `json:"authorId"`
// Author    string `json:"author"`
// Tags      string `json:"tags"`
// Desc      string `json:"desc"`
// Hero      string `json:"hero"`
// Thumb     string `json:"thumb"`
// Body      string `json:"body"`
// Audience  string `json:"audience,omitempty"`
// Weight    int    `json:"weight,omitempty"`
// CreatedOn int64  `json:"createdOn,omitempty"`
// UpdatedOn int64  `json:"updatedOn,omitempty"`
// UpdatedBy string `json:"updatedBy" `

/* QUERIES */

// CreatePost adds a new entry to the database.
func CreatePost(c *gin.Context) {

	db := c.MustGet(model.KeyDb).(*jgorm.DB)
	post := c.MustGet(model.KeyPost).(model.Post)

	if post.Slug != "" && post.Title != "" { // && post.AuthorID != ""
		err := db.Create(&post).Error
		if err != nil {
			msg := fmt.Sprintf("could not create post at: %s, %s", post.Slug, err.Error())
			log.Println(msg)
			c.JSON(503, msg)
			return
		}
		c.JSON(201, gin.H{"success": post})
	} else {
		c.JSON(422, gin.H{"error": "Not enough info"})
	}
}

// DeletePost definitively removes a post from the database.
func DeletePost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*jgorm.DB)

	post := c.MustGet(model.KeyPost).(model.Post)
	slug := c.Params.ByName(model.KeySlug)

	msg := fmt.Sprintf("About to delete post at %s, values: %v", slug, post)
	log.Println(msg)

	// var post model.Post
	d := db.Where(model.KeySlug+" = ?", slug).Delete(&post)
	fmt.Println(d)
	c.JSON(200, gin.H{"Post at  " + slug: "deleted"})
}

// UpdatePost simply updates a post in the database.
func UpdatePost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*jgorm.DB)

	post := c.MustGet(model.KeyPost).(model.Post)
	var toUpdatePost model.Post

	slug := c.Params.ByName(model.KeySlug)
	if err := db.Where(model.KeySlug+" = ?", slug).First(&toUpdatePost).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	// TODO make this cleanly
	toUpdatePost.AuthorID = post.AuthorID
	toUpdatePost.Title = post.Title
	toUpdatePost.Tags = post.Tags
	toUpdatePost.Desc = post.Desc
	toUpdatePost.Hero = post.Hero
	toUpdatePost.Thumb = post.Thumb
	toUpdatePost.Body = post.Body
	toUpdatePost.UpdatedBy = post.UpdatedBy

	if post.PublishedAt != nil {
		toUpdatePost.PublishedAt = post.PublishedAt
	} else {
		toUpdatePost.PublishedAt = &post.CreatedAt
	}

	db.Save(&toUpdatePost)
	c.JSON(200, toUpdatePost)
}

// GetPost simply retrieves a post by slug.
func GetPost(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*jgorm.DB)

	slug := c.Params.ByName(model.KeySlug)
	var post model.Post
	err := db.Where(model.KeySlug+" = ?", slug).First(&post).Error
	if err != nil {
		fmt.Printf("Could not find post with slug %s, cause: %s\n", slug, err.Error())
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, gin.H{"post": post})
	}
}

// GetPosts retrieves all posts, potentially filtered by passed tag value.
func GetPosts(c *gin.Context) {

	queryTag := c.Query(model.KeyTag)
	fmt.Printf("Listing posts for tag [%s]\n", queryTag)
	db := c.MustGet(model.KeyDb).(*jgorm.DB)

	// countLimit := c.Query("count")
	// count, err := strconv.Atoi(countLimit)
	// if err != nil {
	// 	fmt.Printf("Could not deserialize count %s: %s\n", countLimit, err.Error())
	// 	c.Error(err)
	// }
	// fmt.Println("Got a count..." + countLimit)

	var posts []model.Post
	var err error

	if queryTag == "" {
		err = db.Order("updated_at desc").Find(&posts).Error
	} else {
		err = db.Where(model.KeyTags+" LIKE ?", "%"+queryTag+"%").Order("updated_at desc").Find(&posts).Error
		// fmt.Printf("Query: [%v]\n", db.Where("tags LIKE ?", queryTag).Find(&posts).QueryExpr())
	}

	fmt.Printf("Found %d posts for query", len(posts))

	if err != nil {
		c.Error(err)
		c.AbortWithStatus(500)
	} else {
		c.JSON(200, gin.H{"posts": posts})
	}

	// might be enhanced using https://stackoverflow.com/questions/37562873/most-idiomatic-way-to-select-elements-from-an-array-in-golang
	//	var updatedPosts []model.Post
	//	for _, post := range posts {
	//		// if i > count {
	//		// 	break
	//		// }
	//		query := bson.M{model.KeyParentID: bson.RegEx{post.Path, ""}}
	//		post.CommentCount, err = db.C(model.CommentCollection).Find(query).Count()
	//		updatedPosts = append(updatedPosts, post)
	//		fmt.Printf("retrieving post at %s with %d comments \n", post.Path, post.CommentCount)
	//		// TODO also filter by rights
	//	}

}
