package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/model"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	jwtToTest = ""
)

func TestJwtToken(t *testing.T) {
	if jwtToTest != "" {
		err := CheckCredentialAgainstFireBase(nil, jwtToTest)
		if err != nil {
			t.Error("invalid JWT", err.Error())
		}
	}
}

func TestUpdateUsersFromFirebase(t *testing.T) {

	Convey("Test retrieval and update from firebase: ", t, func() {
		// initialise context
		adminEmail = "bruno@sinou.org"
		anonymousEmail = "guest@sinou.org"
		fname := fmt.Sprintf("/sqlite-%d.db", time.Now().Unix())
		path := os.TempDir() + fname
		gorm.InitGormTestRepo(path)
		defer func() {
			os.Remove(path)
		}()

		err := ListExistingUsers(nil)
		if err != nil {
			t.Error("cannot list users", err.Error())
		}

		db := gorm.GetConnection()
		defer db.Close()

		var users []model.User
		db.Preload("Roles").Find(&users)
		for _, currUser := range users {
			printUser(&currUser)
		}
	})
}

func printUser(user *model.User) {
	tstStr, err := json.Marshal(user)
	So(err, ShouldBeNil)
	fmt.Println(string(tstStr))
}
