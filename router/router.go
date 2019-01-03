package router

import (
	"tpark_db/handlers"

	"github.com/buaazp/fasthttprouter"
)

// NewRouter registers URL and handlers then returns instance of router.
func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/api/forum/:slug", handlers.ForumCreateHandler)
	router.POST("/api/forum/:slug/create", handlers.ForumCreateThreadHandler)
	router.GET("/api/forum/:slug/details", handlers.ForumGetOneHandler)
	router.GET("/api/forum/:slug/threads", handlers.ForumGetThreadsHandler)
	router.GET("/api/forum/:slug/users", handlers.ForumGetUsersHandler)

	router.GET("/api/post/:id/details", handlers.PostGetOneHandler)
	router.POST("/api/post/:id/details", handlers.PostUpdateHandler)

	router.GET("/api/service/status", handlers.StatusHandler)
	router.POST("/api/service/clear", handlers.ClearHandler)

	router.GET("/api/thread/:slug_or_id/details", handlers.ThreadGetOneHandler)
	router.GET("/api/thread/:slug_or_id/posts", handlers.ThreadGetPostsHandler)
	router.POST("/api/thread/:slug_or_id/create", handlers.ThreadCreateHandler)
	router.POST("/api/thread/:slug_or_id/details", handlers.ThreadUpdateHandler)
	router.POST("/api/thread/:slug_or_id/vote", handlers.ThreadVoteHandler)

	router.GET("/api/user/:nickname/profile", handlers.UserGetOneHandler)
	router.POST("/api/user/:nickname/create", handlers.UserCreateHandler)
	router.POST("/api/user/:nickname/profile", handlers.UserUpdateHandler)

	return router
}
