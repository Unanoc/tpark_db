package router

import (
	"tpark_db/api/forum"
	"tpark_db/api/post"
	"tpark_db/api/service"
	"tpark_db/api/thread"
	"tpark_db/api/user"

	"github.com/buaazp/fasthttprouter"
)

func NewRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/api/forum/:slug", forum.ForumCreateHandler) //done
	router.POST("/api/forum/:slug/create", forum.ForumCreateThreadHandler) //TODO доделать Helper
	router.GET("/api/forum/:slug/details", forum.ForumGetOneHandler) //done
	router.GET("/api/forum/:slug/threads", forum.ForumGetThreadsHandler) //done
	router.GET("/api/forum/:slug/users", forum.ForumGetUsersHandler) //TODO

	router.GET("/api/post/:id/details", post.PostGetOneHandler)
	router.POST("/api/post/:id/details", post.PostUpdateHandler)

	router.GET("/api/service/status", service.StatusHandler)
	router.POST("/api/service/clear", service.ClearHandler)

	router.GET("/api/thread/:slug_or_id/details", thread.ThreadGetOneHandler)
	router.GET("/api/thread/:slug_or_id/posts", thread.ThreadGetPosts)
	router.POST("/api/thread/:slug_or_id/create", thread.ThreadCreateHandler)
	router.POST("/api/thread/:slug_or_id/details", thread.ThreadUpdateHandler)
	router.POST("/api/thread/:slug_or_id/vote", thread.ThreadVoteHandler)

	router.GET("/api/user/:nickname/profile", user.UserGetOneHandler)
	router.POST("/api/user/:nickname/create", user.UserCreateHandler)
	router.POST("/api/user/:nickname/profile", user.UserUpdateHandler)

	return router
}
