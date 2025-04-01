package server

import (
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (ms *my_server) registerRoutes() http.Handler {
	router := chi.NewRouter()

	h := initWS()
	go h.run()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Get("/auth/{provider}", handleBeginAuth)
	router.Get("/auth/callback/{provider}", ms.handleAuthCallback)
	router.Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello Seattle"))
		})

	v1Router := chi.NewMux()
	v1Router.Use(getUserSession)

	v1Router.Get("/user", ms.getUserMiddleware(handleGetUser))

	// workspace
	v1Router.Get("/workspaces", ms.getUserMiddleware(ms.handleGetUserWorkspaces))
	v1Router.Post("/workspace/create", ms.getUserMiddleware(ms.handleCreateWorkspace))
	v1Router.Delete("/workspace/delete/{workspace_id}", ms.getUserMiddleware(ms.handleDeleteWorkspace))

	// pages
	v1Router.Post("/page/create", ms.handleCreatePage)
	v1Router.Get("/pages/{workspace_id}", ms.getUserMiddleware(ms.handleGetWorkspacePages))
	v1Router.Delete("/page/delete/{page_id}/{workspace_id}", ms.getUserMiddleware(ms.handleDeletePage))

	v1Router.Get("/ws", ms.getUserMiddleware(func(w http.ResponseWriter, r *http.Request, user database.User) {
		new_wsClient(h, w, r, user)
	}))

	router.Mount("/api/v1", v1Router)

	return router
}
