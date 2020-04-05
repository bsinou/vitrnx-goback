package route

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/handler"
	"github.com/bsinou/vitrnx-goback/model"
)

// StartRouter starts the main gin gonic router after configuration.
func StartRouter() {
	r := gin.Default()
	r.Use(loggingHandler(), cors())
	declareRoutes(r)
	log.Fatal(r.Run(":" + conf.GetPort()))
}

func declareRoutes(r *gin.Engine) {

	// API ENTRY POINT
	// Must be logged in
	apiG := r.Group(model.ApiPrefix)
	apiG.Use(verifyToken(), connect())

	// Posts
	posts := apiG.Group("/posts")
	{
		// Configure wrappers for this group
		// posts.Use(addUserMeta(), unmarshallPost(), applyPostPolicies())
		posts.Use(unmarshallPost())

		// Enable fetch with js and CORS
		posts.OPTIONS("", handler.DoNothing)                            // POST
		posts.OPTIONS(":"+model.KeySlug, handler.DoNothing)             // PUT, DELETE
		posts.OPTIONS(":"+model.KeySlug+"/comments", handler.DoNothing) // PUT, DELETE

		// REST
		posts.GET("", handler.GetPosts)                     // query with params
		posts.GET(":"+model.KeySlug, handler.GetPost)       // get one
		posts.POST("", handler.CreatePost)                  // new post
		posts.POST(":"+model.KeySlug, handler.UpdatePost)   // update post
		posts.DELETE(":"+model.KeySlug, handler.DeletePost) // delete post
	}

	// PUBLIC ENTRY POINT
	// Anonymous users can only see public posts, static pages and a few utils pages
	pubG := r.Group(model.PublicPrefix)

	// TODO add limited credentials to also track anonymous user to prevent DDOS and other attacks
	pubG.Use(pubconnect())

	// Posts
	pposts := pubG.Group("/posts")
	{
		// Configure wrappers for this group
		// posts.Use(loggingHandler(), cors(), checkCredentials(), Connect(), addUserMeta(), unmarshallPost(), applyPostPolicies())
		pposts.Use(unmarshallPost())

		// REST
		pposts.GET("", handler.GetPosts)               // query with params
		pposts.GET(":"+model.KeySlug, handler.GetPost) // get one
		// pposts.GET(":"+model.KeySlug+"/comments", handler.ListPostComments) // get post comments
	}
}
