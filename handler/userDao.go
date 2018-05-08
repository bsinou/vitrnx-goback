package handler

// Thanks @etiennerouzeaud to https://gist.github.com/EtienneR/ed522e3d31bc69a9dec3335e639fcf60 && https://medium.com/@etiennerouzeaud/how-to-create-a-basic-restful-api-in-go-c8e032ba3181

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/model"

	// Use SQLite
	_ "github.com/mattn/go-sqlite3"
)

/* QUERIES */

func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []model.User
	db.Find(&users)
	c.JSON(200, users)
}

/* CRUD */

func PutUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*gorm.DB)
	user := c.MustGet(model.KeyUser).(model.User)

	if user.Name != "" && user.Email != "" {
		db.Create(&user)
		c.JSON(201, gin.H{"success": user})
	} else {
		c.JSON(422, gin.H{"error": "Not enough info"})
	}
}

func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Params.ByName("id")
	var user model.User
	db.First(&user, id)

	if user.ID != 0 {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Params.ByName("id")
	var user model.User
	db.First(&user, id)

	// Performs validation
	if user.Name != "" && user.Email != "" {

		if user.ID != 0 {
			var newUser model.User
			c.Bind(&newUser)

			result := model.User{
				// model.ID:        user.ID,
				Name:  newUser.Name,
				Email: newUser.Email,
			}

			db.Save(&result)
			c.JSON(200, gin.H{"success": result})
		} else {
			c.JSON(404, gin.H{"error": "User not found"})
		}

	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get id user
	id := c.Params.ByName("id")
	var user model.User
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.ID != 0 {
		db.Delete(&user)
		c.JSON(200, gin.H{"success": "User #" + id + " deleted"})
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}
