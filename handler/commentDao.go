package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
)

// ListComments retrieves all comments, potentially filtered by the passed parentId.
func ListComments(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	comments := []model.Comment{}
	var err error

	pID := c.Query(model.KeyParentID)
	if pID == "" {
		err = db.C(model.CommentCollection).Find(nil).Sort("-date").All(&comments)
	} else {
		query := bson.M{model.KeyParentID: bson.RegEx{pID, ""}}
		err = db.C(model.CommentCollection).Find(query).Sort("-date").All(&comments)
	}
	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{"comments": comments})
}

/* CRUD */

// PutComment simply creates or updates a comment in the document repository.
func PutComment(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)
	comment := c.MustGet(model.KeyComment).(model.Comment)

	comments := db.C(model.CommentCollection)
	creation := false

	if comment.ID.Hex() == "" {
		creation = true
		// Set creation info
		comment.ID = bson.NewObjectId()
		comment.Date = time.Now().Unix()
		comment.AuthorID = c.MustGet(model.KeyUserID).(string)
		comment.Author = c.MustGet(model.KeyUserDisplayName).(string)
		comment.CreatedOn = time.Now().Unix()
		// } else {
		// 	// TODO clean this: date is lost on update
		// 	oldComment, err := findCommentByID(db, comment.ID)
		// 	if err != nil {
		// 		c.Error(err)
		// 		return
		// 	}
		// 	// fmt.Println("### Updating: received date: " + comment.Date.Format("2006-01-02"))
		// 	// fmt.Println("### Updating: old date: " + oldComment.Date.Format("2006-01-02"))
		// 	comment.Date = oldComment.Date
		// 	comment.AuthorID = oldComment.AuthorID
		// 	comment.Author = oldComment.Author
	}

	// Always update the update (...) info
	comment.UpdatedOn = time.Now().Unix()
	comment.UpdatedBy = c.MustGet(model.KeyUserID).(string)

	if creation {
		err := comments.Insert(comment)
		if err != nil {
			fmt.Printf("Insert failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(201, gin.H{"comment": comment})
	} else {
		query := bson.M{"id": bson.ObjectIdHex(comment.ID.Hex())}
		err := comments.Update(query, comment)
		if err != nil {
			fmt.Printf("Update failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(200, gin.H{"comment": comment})
	}
}

// DeleteComment definitively removes a comment from the repository
func DeleteComment(c *gin.Context) {
	db := c.MustGet(model.KeyDb).(*mgo.Database)

	cID := c.Param(model.KeyMgoID)
	query := bson.M{model.KeyMgoID: bson.ObjectIdHex(cID)}
	err := db.C(model.CommentCollection).Remove(query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Comment has been removed")})
}

/* Helper functions */

func findCommentByID(db *mgo.Database, id bson.ObjectId) (*model.Comment, error) {
	var comment model.Comment
	query := bson.M{"id": bson.ObjectIdHex(id.Hex())}
	err := db.C(model.CommentCollection).Find(query).One(&comment)
	if err != nil {
		err = fmt.Errorf("update failed: could not find comment with id %s. Cause: %v", id, err)
		fmt.Println(err.Error())
		return nil, err
	}
	return &comment, nil
}
