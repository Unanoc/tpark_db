package router

import (
	"fmt"
	"net/http"
	"strings"

	"tpark_db/api"
	"tpark_db/logger"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},

	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		api.ClearHandler,
	},

	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/create",
		api.ForumCreateHandler,
	},

	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/details",
		api.ForumGetOneHandler,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/threads",
		api.ForumGetThreadsHandler,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/users",
		api.ForumGetUsersHandler,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/{id}/details",
		api.PostGetOneHandler,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/{id}/details",
		api.PostUpdateHandler,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/create",
		api.PostsCreateHandler,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		api.StatusHandler,
	},

	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/{slug}/create",
		api.ThreadCreateHandler,
	},

	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/details",
		api.ThreadGetOneHandler,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/posts",
		api.ThreadGetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/details",
		api.ThreadUpdateHandler,
	},
	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/vote",
		api.ThreadVoteHandler,
	},

	Route{
		"UserCreate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/create",
		api.UserCreateHandler,
	},

	Route{
		"UserGetOne",
		strings.ToUpper("Get"),
		"/api/user/{nickname}/profile",
		api.UserGetOneHandler,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/profile",
		api.UserUpdateHandler,
	},
}
