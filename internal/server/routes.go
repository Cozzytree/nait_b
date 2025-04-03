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
	v1Router.Get("/workspaces",
		ms.getUserMiddleware(ms.handleGetUserWorkspaces))
	v1Router.Get("/workspace/members/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceMembers)))

	v1Router.Post("/workspace/create",
		ms.getUserMiddleware(ms.handleCreateWorkspace))
	v1Router.Delete("/workspace/delete/{workspace_id}",
		ms.getUserMiddleware(ms.handleDeleteWorkspace))

	// pages
	v1Router.Post("/page/create", ms.handleCreatePage)
	v1Router.Get("/pages/{workspace_id}",
		ms.getUserMiddleware(
			ms.handleGetWorkspacePages))
	v1Router.Delete("/page/delete/{page_id}/{workspace_id}",
		ms.getUserMiddleware(
			ms.handleDeletePage))

	// tasks
	v1Router.Post("/task/create",
		ms.getUserMiddleware(
			ms.handleCreateTask))

	v1Router.Get("/task/count/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetTotalTaskCount)))
	v1Router.Get("/task/s_count/{status}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceTaskStatusCount)))
	v1Router.Get("/task/p_count/{priority}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceTaskPriorityCount)))
	v1Router.Get("/task/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceTasks)))
	v1Router.Get("/task/completed/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceCompletedTask)))
	v1Router.Get("/task/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetTask)))
	v1Router.Get("/task/childs/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetChildTasks)))
	v1Router.Get("/task/assigned/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceUsersAssignedTask)))
	v1Router.Get("/task/created/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceUsersCreatedTask)))

	v1Router.Patch("/task/desc/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskDesciption)))
	v1Router.Patch("/task/priority/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskPriority)))
	v1Router.Patch("/task/status/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskStatus)))
	v1Router.Patch("/task/name/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskName)))
	v1Router.Patch("/task/due/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskDue)))

	v1Router.Delete("/task/delete/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleDeleteTask)))

	v1Router.Get("/ws", ms.getUserMiddleware(func(w http.ResponseWriter, r *http.Request, user database.User) {
		new_wsClient(h, w, r, user)
	}))

	router.Mount("/api/v1", v1Router)

	return router
}
