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
	// r.Use(loggingHandler(), cors(), checkCredentials(), applyPolicies())
	declareRoutes(r)
	r.Run(":8888")
}

func declareRoutes(r *gin.Engine) {

	api := r.Group("/api")

	// Users
	user := api.Group("/users")
	{
		t := model.StoreTypeGorm
		// Enable fetch with js and CORS
		user.OPTIONS("/", optionsUser)    // POST
		user.OPTIONS("/:id", optionsUser) // PUT, DELETE
		// REST
		user.POST("/", Connect(t), handler.PostUser)
		user.GET("/", Connect(t), handler.GetUsers)
		user.GET("/:id", Connect(t), handler.GetUser)
		user.PUT("/:id", Connect(t), handler.UpdateUser)
		user.DELETE("/:id", Connect(t), handler.DeleteUser)

		// // shortcut to backend type
		// t := model.StoreTypeGorm
		// // Enable fetch with js and CORS
		// api.OPTIONS("/users", optionsUser)     // POST
		// api.OPTIONS("/users/:id", optionsUser) // PUT, DELETE
		// // REST
		// api.POST("/users", Connect(t), handler.PostUser)
		// api.GET("/users", Connect(t), handler.GetUsers)
		// api.GET("/users/:id", Connect(t), handler.GetUser)
		// api.PUT("/users/:id", Connect(t), handler.UpdateUser)
		// api.DELETE("/users/:id", Connect(t), handler.DeleteUser)
	}

	// Posts
	posts := r.Group("/api/posts")
	{
		// shortcut to backend type
		t := model.StoreTypeMgo
		// Configure wrappers for this group
		posts.Use(loggingHandler(), cors(), checkCredentials(), applyPolicies())

		// Enable fetch with js and CORS
		posts.OPTIONS("", optionsUser)    // POST
		posts.OPTIONS(":id", optionsUser) // PUT, DELETE

		// REST
		posts.GET("", Connect(t), handler.ListPosts)                    // query with params
		posts.GET(":"+model.KeyPath, Connect(t), handler.ReadPost)      // get one
		posts.POST("", Connect(t), handler.PutPost)                     // new post
		posts.POST(":"+model.KeyPath, Connect(t), handler.PutPost)      // update post
		posts.DELETE(":"+model.KeyPath, Connect(t), handler.DeletePost) // delete post
	}
}

func doNothing(c *gin.Context) {
	fmt.Println("We should never reach this point, but received a request at " + c.Request.URL.String())
	c.Next()
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
