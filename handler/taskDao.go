package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
)

func ListTasks(c *gin.Context) {

	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	userID := c.MustGet(model.KeyUserID).(string)
	roles := c.MustGet(model.KeyUserRoles).([]string)

	tasks := []model.Task{}
	var err error

	catID := c.Query(model.KeyCategoryID)

	showAll := c.Query("showAll")
	showClosed := c.Query("showClosed")

	query := bson.M{}
	if !(catID == "" || catID == "all") {
		query[model.KeyCategoryID] = bson.RegEx{catID, ""}
	}

	if showAll != "true" {
		query["managerId"] = map[string]interface{}{"$in": append(roles, userID)}
	}

	if !(showClosed == "true") {
		query["closeDate"] = 0
	}

	err = db.C(model.TaskCollection).Find(query).Sort("-flags", "-dueDate").All(&tasks)
	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{"tasks": tasks})
}

/* CRUD */

func PutTask(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)
	task := c.MustGet(model.KeyTask).(model.Task)

	tasks := db.C(model.TaskCollection)
	creation := false

	if task.ID.Hex() == "" {
		creation = true
		// Set creation info
		task.ID = bson.NewObjectId()
		task.CreationDate = time.Now().Unix()
		task.ManagerID = c.MustGet(model.KeyUserID).(string)
		task.Manager = c.MustGet(model.KeyUserDisplayName).(string)
		task.DueDate = time.Now().AddDate(0, 0, 14).Unix()
	}

	// Always update the update (...) info
	task.UpdateDate = time.Now().Unix()
	task.UpdateBy = c.MustGet(model.KeyUserID).(string)

	if creation {
		err := tasks.Insert(task)
		if err != nil {
			fmt.Printf("Insert failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(201, gin.H{"task": task})
	} else {
		query := bson.M{"id": bson.ObjectIdHex(task.ID.Hex())}
		err := tasks.Update(query, task)
		if err != nil {
			fmt.Printf("Update failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(200, gin.H{"task": task})
	}
}

// DeleteTask definitively removes a task from the repository
func DeleteTask(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	tID := c.Param(model.KeyMgoID)
	query := bson.M{model.KeyMgoID: bson.ObjectIdHex(tID)}
	err := db.C(model.TaskCollection).Remove(query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Task has been removed")})
}
