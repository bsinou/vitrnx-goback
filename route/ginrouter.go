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
		// api.GET("/new", Connect(t), handler.NewPost)
		api.POST("/posts", Connect(t), handler.CreatePost)
		api.GET("/posts/:_id", Connect(t), handler.EditPost)
		api.GET("/posts", Connect(t), handler.ListPosts)
		api.POST("/posts/:_id", Connect(t), handler.UpdatePost)
		api.DELETE("/posts/:_id", Connect(t), handler.DeletePost)
	}
}

func optionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}
