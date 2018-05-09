package handler

// Thanks @etiennerouzeaud to https://gist.github.com/EtienneR/ed522e3d31bc69a9dec3335e639fcf60 && https://medium.com/@etiennerouzeaud/how-to-create-a-basic-restful-api-in-go-c8e032ba3181

import (
	"log"

	"github.com/gin-gonic/gin"
	jgorm "github.com/jinzhu/gorm"

	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"

	// Use SQLite
	_ "github.com/mattn/go-sqlite3"
)

/* QUERIES */

func GetUsers(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	var users []model.User
	db.Preload("Roles").Find(&users)

	// TODO manage adding roles in the list
	c.JSON(200, gin.H{"users": users})
}

/* CRUD */

func PutUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	user := c.MustGet(model.KeyUser).(model.User)

	// TODO Check:
	// -> Put simple registered role on creation
	// -> explicitly copy editable properties when editing self
	// -> double check permission when editing roles
	// -> only admin users can change admin & user admin roles

	if user.Name != "" && user.Email != "" {
		db.Create(&user)
		c.JSON(201, gin.H{"success": user})
	} else {
		c.JSON(422, gin.H{"error": "Not enough info"})
	}
}

func GetUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	id := c.Params.ByName("id")
	var user model.User
	// err := db.Preload("Roles").Where(&model.User{UserID: id}).First(&user).Error
	err := db.Preload("Roles").First(&user, id).Error
	if err != nil {
		log.Println("could not retrieve user: " + err.Error())
		c.JSON(503, "User not found, server error")
		return
	}

	gorm.WithUserRoles(&user)

	if user.ID != 0 {
		c.JSON(200, gin.H{"user": user})
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)

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
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)

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
