package router

import (
	"tpark_db/api"
	
	"github.com/buaazp/fasthttprouter"
)

func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/api/forum/:slug/details", api.ForumGetOneHandler)
	router.GET("/api/forum/:slug/threads", api.ForumGetThreadsHandler)
	router.GET("/api/forum/:slug/users", api.ForumGetUsersHandler)
	router.POST("/api/forum/create", api.ForumCreateHandler)	
	router.POST("/api/forum/:slug_or_id/create", api.ThreadCreateHandler)	

	router.GET("/api/post/:id/details", api.PostGetOneHandler)
	router.POST("/api/post/:id/details", api.PostUpdateHandler)	

	router.GET("/api/service/status", api.StatusHandler)
	router.POST("/api/service/clear", api.ClearHandler)

	router.GET("/api/thread/:slug_or_id/details", api.ThreadGetOneHandler)
	router.GET("/api/thread/:slug_or_id/posts", api.ThreadGetPosts)
	router.POST("/api/thread/:slug_or_id/create", api.PostsCreateHandler)	
	router.POST("/api/thread/:slug_or_id/details", api.ThreadUpdateHandler)	
	router.POST("/api/thread/:slug_or_id/vote", api.ThreadVoteHandler)

	router.GET("/api/user/:nickname/profile", api.UserGetOneHandler)
	router.POST("/api/user/:nickname/create", api.UserCreateHandler)	
	router.POST("/api/user/:nickname/profile", api.UserUpdateHandler)		

	return router
}