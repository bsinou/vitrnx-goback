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
		authG.Use(loggingHandler(), cors(), checkCredentials(), Connect())
		authG.OPTIONS("login", handler.DoNothing) // POST
		authG.POST("login", auth.PostLogin)
	}

	// Users
	user := r.Group(model.ApiPrefix + "users")
	{
		// Configure wrappers for this group
		user.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta())
		// Enable fetch with js and CORS
		user.OPTIONS("", handler.DoNothing)
		user.OPTIONS(":id", handler.DoNothing)
		user.OPTIONS(":id/roles", handler.DoNothing)

		// REST
		user.GET("", handler.GetUsers)                                 // query with params
		user.GET(":"+model.KeyUserID, handler.GetUser)                 // get one
		user.POST("", applyUserCreationPolicies(), handler.CreateUser) // CREATION
		user.PATCH(":id", applyUserUpdatePolicies(), handler.PatchUser)
		user.PATCH(":id/roles", applyUserRolesUpdatePolicies(), handler.PatchUserRoles)
	}

	// Roles
	roles := r.Group(model.ApiPrefix + "roles")
	{
		roles.Use(loggingHandler(), cors(), checkCredentials(), Connect())
		roles.OPTIONS("", handler.DoNothing)    // POST
		roles.OPTIONS(":id", handler.DoNothing) // PUT, DELETE

		// REST
		roles.GET("", handler.GetRoles) // query with params
		// user.GET(":"+model.KeyUserID, handler.GetUser) // get one
		// user.POST("", checkCredentialsForUserCreation(), handler.PutUser)
	}

	// Posts
	posts := r.Group(model.ApiPrefix + "posts")
	{
		// Configure wrappers for this group

		// Split JWT check and user meta

		// FIXME
		posts.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta(), unmarshallPost(), applyPostPolicies())
		// posts.Use(loggingHandler(), cors(), Connect(), checkCredentials(), unmarshallPost(), applyPostPolicies())

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
		comments.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta(), unmarshallComment(), applyCommentPolicies())
		comments.OPTIONS("", handler.DoNothing)
		comments.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)

		// REST
		comments.GET("", handler.ListComments)                     // query with params
		comments.POST("", handler.PutComment)                      // Create or update
		comments.DELETE(":"+model.KeyMgoID, handler.DeleteComment) // delete comment
	}
}
