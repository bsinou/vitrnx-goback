package route

import (
	"fmt"

	"github.com/bsinou/vitrnx-goback/handler"
	"github.com/bsinou/vitrnx-goback/model"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Starting Gin router")
	r := gin.Default()
	r.Use(loggingHandler(), cors(), checkCredentials())
	declareRoutes(r)
	r.Run(":8888")
}

func declareRoutes(r *gin.Engine) {

	api := r.Group("api")

	// Users
	{
		// shortcut to backend type
		t := model.StoreTypeGorm
		// Enable fetch with js and CORS
		api.OPTIONS("/users", optionsUser)     // POST
		api.OPTIONS("/users/:id", optionsUser) // PUT, DELETE
		// REST
		api.POST("/users", Connect(t), handler.PostUser)
		api.GET("/users", Connect(t), handler.GetUsers)
		api.GET("/users/:id", Connect(t), handler.GetUser)
		api.PUT("/users/:id", Connect(t), handler.UpdateUser)
		api.DELETE("/users/:id", Connect(t), handler.DeleteUser)
	}

	// Posts
	{
		// shortcut to backend type
		t := model.StoreTypeMgo
		// Enable fetch with js and CORS
		api.OPTIONS("/posts", optionsUser)     // POST
		api.OPTIONS("/posts/:id", optionsUser) // PUT, DELETE
		// REST
		api.GET("/posts/:"+model.KeyPath, Connect(t), handler.ReadPost)      // get one
		api.GET("/posts", Connect(t), handler.ListPosts)                     // query with params
		api.POST("/posts", Connect(t), handler.PutPost)                      // new post
		api.POST("/posts/:"+model.KeyPath, Connect(t), handler.PutPost)      // update post
		api.DELETE("/posts/:"+model.KeyPath, Connect(t), handler.DeletePost) // delete post
	}
}

func optionsUser(c *gin.Context) {
	fmt.Println("Received an OPTIONS request at " + c.Request.URL.String())

	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
	// Rather use this than below lines
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// The second 'Authorization' line erase the first and Content-Type is not an authorized header anymore
	// Thus it's does not work
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")

	c.Next()
}
