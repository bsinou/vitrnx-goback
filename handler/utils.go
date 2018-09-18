package handler

import (
	"fmt"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/gin-gonic/gin"
)

// DoNothing simply forward to next handler
func DoNothing(c *gin.Context) {
	fmt.Println("We should never reach this point, but received a request at " + c.Request.URL.String())
	c.Next()
}

// BasicCheck returns a simple page with some info about the current installation
func BasicCheck(c *gin.Context) {
	c.JSON(200, gin.H{"Instance ID": conf.VitrnxInstanceID, "Version": conf.VitrnxVersion, "Env": conf.Env, "Build timestamp": conf.BuildTimestamp})
}
