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
		authG.Use(loggingHandler(), cors(), Connect(), checkCredentials())
		authG.OPTIONS("login", handler.DoNothing) // POST
		authG.POST("login", auth.PostLogin)
	}

	// Users
	user := r.Group(model.ApiPrefix + "users")
	{
		// Configure wrappers for this group
		user.Use(loggingHandler(), cors(), Connect())
		// Enable fetch with js and CORS
		user.OPTIONS("", handler.DoNothing)    // POST
		user.OPTIONS(":id", handler.DoNothing) // PUT, DELETE

		// REST
		user.GET("", handler.GetUsers) // query with params
		user.POST("", checkCredentialsForUserCreation(), handler.PutUser)
		user.GET(":"+model.KeyUserID, handler.GetUser) // get one
	}

	// Posts
	posts := r.Group(model.ApiPrefix + "posts")
	{
		// Configure wrappers for this group
		posts.Use(loggingHandler(), cors(), Connect(), checkCredentials(), unmarshallPost(), applyPostPolicies())

		// Enable fetch with js and CORS
		posts.OPTIONS("", handler.DoNothing)    // POST
		posts.OPTIONS(":id", handler.DoNothing) // PUT, DELETE

		// REST
		posts.GET("", handler.ListPosts)                    // query with params
		posts.GET(":"+model.KeyPath, handler.ReadPost)      // get one
		posts.POST("", handler.PutPost)                     // new post
		posts.POST(":"+model.KeyPath, handler.PutPost)      // update post
		posts.DELETE(":"+model.KeyPath, handler.DeletePost) // delete post
	}

	// Comments
	comments := r.Group(model.ApiPrefix + "comments")
	{
		comments.Use(loggingHandler(), cors(), Connect(), checkCredentials(), unmarshallComment(), applyCommentPolicies())
		comments.OPTIONS("", handler.DoNothing)
		comments.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)

		// REST
		comments.GET("", handler.ListComments)                     // query with params
		comments.POST("", handler.PutComment)                      // Create or update
		comments.DELETE(":"+model.KeyMgoID, handler.DeleteComment) // delete comment
	}
}
