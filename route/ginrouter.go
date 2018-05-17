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
	r.Use(loggingHandler(), cors(), checkCredentials(), Connect())
	declareRoutes(r)
	log.Fatal(r.Run(":8888"))
}

func declareRoutes(r *gin.Engine) {

	// Authentication
	authG := r.Group(model.ApiPrefix + "auth")
	{
		// authG.Use()
		authG.OPTIONS("login", handler.DoNothing) // POST
		authG.POST("login", auth.PostLogin)
	}

	// Users
	user := r.Group(model.ApiPrefix + "users")
	{
		// Configure wrappers for this group
		// user.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta())
		user.Use(addUserMeta())
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

	// UserMeta
	meta := r.Group(model.ApiPrefix + "usermeta")
	{
		meta.Use(addUserMeta())
		meta.OPTIONS("", handler.DoNothing)
		meta.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)
		// meta.OPTIONS("/by-day", handler.DoNothing)

		// TODO implement
		// meta.GET("/by-day", handler.GetByDay)

		// REST
		meta.GET(":"+model.KeyUserID, handler.ReadPresence)
		meta.POST(":"+model.KeyMgoID, applyUserCreationPolicies(), handler.PutPresence)
	}

	// Roles
	roles := r.Group(model.ApiPrefix + "roles")
	{
		roles.OPTIONS("", handler.DoNothing)                 // POST
		roles.OPTIONS(":"+model.KeyMgoID, handler.DoNothing) // PUT, DELETE
		roles.GET("", handler.GetRoles)                      // query with params
	}

	// Posts
	posts := r.Group(model.ApiPrefix + "posts")
	{
		// Configure wrappers for this group
		// posts.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta(), unmarshallPost(), applyPostPolicies())
		posts.Use(addUserMeta(), unmarshallPost(), applyPostPolicies())

		// Enable fetch with js and CORS
		posts.OPTIONS("", handler.DoNothing)                            // POST
		posts.OPTIONS(":"+model.KeyPath, handler.DoNothing)             // PUT, DELETE
		posts.OPTIONS(":"+model.KeyPath+"/comments", handler.DoNothing) // PUT, DELETE

		// REST
		posts.GET("", handler.ListPosts)                                   // query with params
		posts.GET(":"+model.KeyPath, handler.ReadPost)                     // get one
		posts.GET(":"+model.KeyPath+"/comments", handler.ListPostComments) // get post comments
		posts.POST("", handler.PutPost)                                    // new post
		posts.POST(":"+model.KeyPath, handler.PutPost)                     // update post
		posts.DELETE(":"+model.KeyPath, handler.DeletePost)                // delete post
	}

	// Comments
	comments := r.Group(model.ApiPrefix + "comments")
	{
		// comments.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta(), unmarshallComment(), applyCommentPolicies())
		comments.Use(addUserMeta(), unmarshallComment(), applyCommentPolicies())
		comments.OPTIONS("", handler.DoNothing)
		comments.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)

		// REST
		comments.GET("", handler.ListComments)                     // query with params
		comments.POST("", handler.PutComment)                      // Create or update
		comments.DELETE(":"+model.KeyMgoID, handler.DeleteComment) // delete comment
	}

	// Tasks
	tasks := r.Group(model.ApiPrefix + "tasks")
	{
		tasks.Use(addUserMeta(), unmarshallTask(), applyTaskPolicies())
		tasks.OPTIONS("", handler.DoNothing)
		// tasks.OPTIONS(":"+model.KeyCategoryID, handler.DoNothing)
		tasks.OPTIONS(":"+model.KeyMgoID, handler.DoNothing)

		// REST
		tasks.GET("", handler.ListTasks)
		tasks.GET(":"+model.KeyCategoryID, handler.ListTasks)
		tasks.POST("", handler.PutTask)
		tasks.POST(":"+model.KeyMgoID, handler.PutTask)
		tasks.DELETE(":"+model.KeyMgoID, handler.DeleteTask)
	}

}
