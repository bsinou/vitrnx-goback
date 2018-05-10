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
	err := db.Preload("Roles").Find(&users).Error
	if err != nil {
		log.Println("could not list users: " + err.Error())
		c.JSON(503, "cannot list users")
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func GetRoles(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	var roles []model.Role
	err := db.Find(&roles).Error

	if err != nil {
		log.Println("could not list roles: " + err.Error())
		c.JSON(503, "cannot list roles")
		return
	}

	// TODO manage adding roles in the list
	c.JSON(200, gin.H{"roles": roles})
}

/* CRUD */

func CreateUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	user := c.MustGet(model.KeyEditedUser).(model.User)

	var dflt model.Role
	db.Where(&model.Role{RoleID: model.RoleRegistered}).First(&dflt)
	user.Roles = []model.Role{dflt}

	if user.Name != "" && user.Email != "" {
		err := db.Create(&user).Error
		if err != nil {
			msg := "could not create user: " + err.Error()
			log.Println(msg)
			c.JSON(503, msg)
			return
		}
		c.JSON(201, gin.H{"success": user})
	} else {
		c.JSON(422, gin.H{"error": "Not enough info"})
	}
}

func PatchUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	editedUser := c.MustGet(model.KeyEditedUser).(model.User)

	var loadedUser model.User
	err := db.Where(&model.User{UserID: editedUser.UserID}).First(&loadedUser).Error
	if err != nil {
		log.Println("could not retrieve user: " + err.Error())
		c.JSON(503, "User not found, server error")
		return
	}

	// Add an indirection to prevent updating reserved info
	if editedUser.Name != "" {
		loadedUser.Name = editedUser.Name
	}
	if editedUser.Email != "" {
		loadedUser.Email = editedUser.Email
	}
	if editedUser.Address != "" {
		loadedUser.Address = editedUser.Address
	}

	err = db.Save(&loadedUser).Error
	if err != nil {
		log.Println("could not update user: " + err.Error())
		c.JSON(503, "Server error while updating user")
		return
	}

	c.JSON(200, gin.H{"user": loadedUser})
}

func PatchUserRoles(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	toEditUserID := c.Param("id")
	retrievedRoles := c.MustGet(model.KeyEditedUserRoles).([]string)

	var loadedUser model.User
	err := db.Preload("Roles").Where(&model.User{UserID: toEditUserID}).First(&loadedUser).Error
	if err != nil {
		log.Println("could not retrieve user: " + err.Error())
		c.JSON(503, "User not found, server error")
		return
	}

	roles := make([]*model.Role, len(retrievedRoles))
	for i, r := range retrievedRoles {
		var ro model.Role
		db.Where(&model.Role{RoleID: r}).First(&ro)
		roles[i] = &ro
	}

	err = db.Model(&loadedUser).Association("Roles").Replace(roles).Error
	if err != nil {
		msg := "could not update user roles: " + err.Error()
		log.Println()
		c.JSON(503, msg)
		return
	}

	c.JSON(200, gin.H{"user": loadedUser})
}

func GetUser(c *gin.Context) {
	db := c.MustGet(model.KeyUserDb).(*jgorm.DB)
	id := c.Params.ByName(model.KeyUserID)
	var user model.User
	err := db.Preload("Roles").Where(&model.User{UserID: id}).First(&user).Error
	// err := db.Preload("Roles").First(&user, id).Error
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

	id := c.Params.ByName(model.KeyUserID)
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
	id := c.Params.ByName(model.KeyUserID)
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
