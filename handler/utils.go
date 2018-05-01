package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func DoNothing(c *gin.Context) {
	fmt.Println("We should never reach this point, but received a request at " + c.Request.URL.String())
	c.Next()
}
