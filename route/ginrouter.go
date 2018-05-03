package route

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/handler"
	"github.com/bsinou/vitrnx-goback/model"
)

// func init() {
// 	fmt.Println("Starting Gin router")
// 	r := gin.Default()
// 	// r.Use(loggingHandler(), cors(), checkCredentials(), applyPolicies())
// 	declareRoutes(r)
// 	go r.Run(":8888")
// }

func StartRouter() {
	r := gin.Default()
	// r.Use(loggingHandler(), cors(), checkCredentials(), applyPolicies())
	declareRoutes(r)
	log.Fatal(r.Run(":8888"))
}

func declareRoutes(r *gin.Engine) {

	// Users
	user := r.Group("/api/users")
	{
		// shortcut to backend type
		t := model.StoreTypeGorm
		// Configure wrappers for this group
		user.Use(loggingHandler(), cors())
		// Enable fetch with js and CORS
		user.OPTIONS("", handler.DoNothing)    // POST
		user.OPTIONS(":id", handler.DoNothing) // PUT, DELETE

		// REST
		user.POST("", Connect(t), checkCredentialsForUserCreation(), handler.PostUser)
		// user.GET("", Connect(t), handler.GetUsers)
		// user.GET(":id", Connect(t), handler.GetUser)
		// user.PUT(":id", Connect(t), handler.UpdateUser)
		// user.DELETE(":id", Connect(t), handler.DeleteUser)
	}

	// Posts
	posts := r.Group("/api/posts")
	{
		// shortcut to backend type
		t := model.StoreTypeMgo
		// Configure wrappers for this group
		posts.Use(loggingHandler(), cors(), checkCredentials(), applyPostPolicies())

		// Enable fetch with js and CORS
		posts.OPTIONS("", handler.DoNothing)    // POST
		posts.OPTIONS(":id", handler.DoNothing) // PUT, DELETE

		// REST
		posts.GET("", Connect(t), handler.ListPosts)                    // query with params
		posts.GET(":"+model.KeyPath, Connect(t), handler.ReadPost)      // get one
		posts.POST("", Connect(t), handler.PutPost)                     // new post
		posts.POST(":"+model.KeyPath, Connect(t), handler.PutPost)      // update post
		posts.DELETE(":"+model.KeyPath, Connect(t), handler.DeletePost) // delete post
	}
}

// func optionsUser(c *gin.Context, allowedMethods string) {
// 	fmt.Println("Received an OPTIONS request at " + c.Request.URL.String())

// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
// 	// Rather use this than below lines
// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	// The second 'Authorization' line erase the first and Content-Type is not an authorized header anymore
// 	// Thus it's does not work
// 	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")

// 	c.Next()
// }
