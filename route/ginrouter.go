package route

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/auth"
	"github.com/bsinou/vitrnx-goback/handler"
	"github.com/bsinou/vitrnx-goback/model"
)

func StartRouter() {
	r := gin.Default()
	declareRoutes(r)
	log.Fatal(r.Run(":8888"))
}

func declareRoutes(r *gin.Engine) {

	// Authentication
	authG := r.Group(model.ApiPrefix + "auth")
	{
		// shortcut to backend type
		t := model.StoreTypeGorm

		authG.Use(loggingHandler(), cors(), checkCredentials())
		authG.OPTIONS("login", handler.DoNothing) // POST
		authG.POST("login", Connect(t), auth.PostLogin)

	}

	// Users
	user := r.Group(model.ApiPrefix + "users")
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
	posts := r.Group(model.ApiPrefix + "posts")
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

	// Comments
	comments := r.Group(model.ApiPrefix + "comments")
	{
		// shortcut to backend type
		t := model.StoreTypeMgo
		comments.Use(loggingHandler(), cors(), checkCredentials())
		comments.OPTIONS("", handler.DoNothing)
		comments.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)

		// REST
		comments.GET("", Connect(t), handler.ListComments)                     // query with params
		comments.POST("", Connect(t), handler.PutComment)                      // Create or update
		comments.DELETE(":"+model.KeyMgoID, Connect(t), handler.DeleteComment) // delete comment
	}
}
