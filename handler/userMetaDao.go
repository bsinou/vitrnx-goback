package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/bsinou/vitrnx-goback/model"
)

// ListPresences retrieves the presence list
func ListPresences(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	presences := []model.Presence{}

	err := db.C(model.CommentCollection).Find(nil).All(&presences)
	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{"presences": presences})
}

type byDay struct {
	ID      string
	Label   string
	AdultNb int
	ChildNb int
}

// var days []string = []string{"D1","D2","D3","D4","D5","D6","D7","D8"}

func GuestTotal(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	presences := []model.Presence{}
	fmt.Println("Here")

	err := db.C(model.PresenceCollection).Find(nil).All(&presences)
	if err != nil {
		c.Error(err)
	}

	totalAdults := 0
	totalChildren := 0
	for _, p := range presences {
		if !p.IsComing {
			continue
		}
		totalAdults += int(p.AdultNb)
		totalChildren += int(p.ChildNb)
	}

	// c.JSON(200, gin.H{"guestsByDay": map[string]string{
	// 	"totalAdults":   strconv.Itoa(totalChildren),
	// 	"totalChildren": strconv.Itoa(totalChildren),
	// }})
	c.JSON(200, gin.H{"guestsByDay": map[string]int{
		"totalAdults":   totalAdults,
		"totalChildren": totalChildren,
	}})

}

func ListGuestsByDay(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	presences := []model.Presence{}

	err := db.C(model.PresenceCollection).Find(nil).All(&presences)
	if err != nil {
		c.Error(err)
	}

	// well well well ...
	var res [8]byDay
	for _, p := range presences {
		if !p.IsComing {
			continue
		}
		if p.D1 {
			res[0].AdultNb += int(p.AdultNb)
			res[0].ChildNb += int(p.ChildNb)
		}
		if p.D2 {
			res[1].AdultNb += int(p.AdultNb)
			res[1].ChildNb += int(p.ChildNb)
		}
		if p.D3 {
			res[2].AdultNb += int(p.AdultNb)
			res[2].ChildNb += int(p.ChildNb)
		}
		if p.D4 {
			res[3].AdultNb += int(p.AdultNb)
			res[3].ChildNb += int(p.ChildNb)
		}
		if p.D5 {
			res[4].AdultNb += int(p.AdultNb)
			res[4].ChildNb += int(p.ChildNb)
		}
		if p.D6 {
			res[5].AdultNb += int(p.AdultNb)
			res[5].ChildNb += int(p.ChildNb)
		}
		if p.D7 {
			res[6].AdultNb += int(p.AdultNb)
			res[6].ChildNb += int(p.ChildNb)
		}
		if p.D8 {
			res[7].AdultNb += int(p.AdultNb)
			res[7].ChildNb += int(p.ChildNb)
		}
	}

	c.JSON(200, gin.H{"guestsByDay": res})
}

/* CRUD */

// PutPresence simply creates or updates a guest presence in the document repository.
func PutPresence(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	presence := model.Presence{}
	err := c.Bind(&presence)
	if err != nil {
		fmt.Printf("Could not bind presence %v\n", err)
		c.Error(err)
		c.Abort()
		return
	}

	presences := db.C(model.PresenceCollection)
	creation := false

	if presence.ID.Hex() == "" {
		creation = true
		// Set creation info
		presence.ID = bson.NewObjectId()

		// if presence.UserID == "" {
		// 	presence.UserID = c.MustGet(model.KeyUserID).(string)
		// }

		// else {
		// 	var existing model.Presence
		// 	query := bson.M{model.KeyUserID: presence.UserID}
		// 	err := db.C(model.PresenceCollection).Find(query).One(&existing)
		// 	if err != nil {
		// 		fmt.Printf("Insert failed: %s\n", err.Error())
		// 		c.Error(err)
		// 		return
		// 	}
		// 	if presence.ID.Hex() != "" {
		// 		err2 := fmt.Errorf("already existing, cannot create") // already existing, cannot create
		// 		fmt.Println(err2.Error())
		// 		c.Error(err2)
		// 		return
		// 	}
		// }
	}

	// Always update the update (...) info
	presence.UpdatedOn = time.Now().Unix()
	presence.UpdatedBy = c.MustGet(model.KeyUserID).(string)

	if creation {
		err := presences.Insert(presence)
		if err != nil {
			fmt.Printf("Insert failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(201, gin.H{"presence": presence})
	} else {
		query := bson.M{"id": bson.ObjectIdHex(presence.ID.Hex())}
		err := presences.Update(query, presence)
		if err != nil {
			fmt.Printf("Update failed: %s\n", err.Error())
			c.Error(err)
		}
		c.JSON(200, gin.H{"presence": presence})
	}
}

func ReadPresence(c *gin.Context) {
	db := c.MustGet(model.KeyDataDb).(*mgo.Database)

	userID := c.Param(model.KeyUserID)
	fmt.Printf("Reading presence for: %s\n", userID)
	query := bson.M{model.KeyUserID: userID}
	presence := model.Presence{}
	err := db.C(model.PresenceCollection).Find(query).One(&presence)
	if err != nil {
		fmt.Printf("Could not read presence of %s: %s\n", c.Param(model.KeyUserID), err.Error())
		c.Error(err)
		return
	}
	// TODO check if current user can see this
	c.JSON(201, gin.H{"presence": presence})
}
