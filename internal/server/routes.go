package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

func (ms *my_server) registerRoutes() http.Handler {
	router := chi.NewRouter()
	limiter := newLimiter()

	h := initWS()
	go h.run()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(gracefulerror)
	router.Use(limiter.limiter)

	router.Get("/auth/{provider}", handleBeginAuth)
	router.Get("/auth/callback/{provider}", ms.handleAuthCallback)
	router.Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello Seattle"))
		})

	router.Get("/invite/{link_id}", func(w http.ResponseWriter, r *http.Request) {
		session, err := gothic.Store.Get(r, os.Getenv("SESSION_NAME"))
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error/msg=%v", os.Getenv("REDIRECT_URL"), "no-session"), http.StatusPermanentRedirect)
			ResponseWithError(w, err.Error(), 401)
			return
		}
		userId, ok := session.Values["user_id"]
		if !ok {
			ResponseWithError(w, "invalid user id", 401)
			return
		}

		auth_id, ok := userId.(string)
		if !ok {
			ResponseWithError(w, "invalid user id", 401)
			return
		}

		user, err := ms.db.GetUser(r.Context(), auth_id)
		if err != nil {
			ResponseWithError(w, err.Error(), 401)
			return
		}

		link_id, err := uuid.Parse(r.PathValue("link_id"))
		if err != nil {
			ResponseWithError(w, "invalid link", 400)
			return
		}

		link, err := ms.db.GetAlink(r.Context(), link_id)
		if err != nil {
			if err == sql.ErrNoRows {
				ResponseWithError(w, "link not found", 404)
			} else {
				ResponseWithError(w, err.Error(), 500)
			}
			return
		}

		_, err = ms.db.GetWorkspaceUserRole(r.Context(), database.GetWorkspaceUserRoleParams{
			UserID:      user.ID,
			WorkspaceID: link.WorkspaceID,
		})
		if err == nil {
			ResponseWithError(w, "already a member", 400)
			return
		}

		err = ms.db.CreateNewWorkspaceMember(r.Context(), database.CreateNewWorkspaceMemberParams{
			UserID:      user.ID,
			WorkspaceID: link.WorkspaceID,
			Role:        link.Role,
		})
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error/msg=%v",
				os.Getenv("REDIRECT_URL"), err.Error()), http.StatusPermanentRedirect)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%s/%s", os.Getenv("REDIRECT_URL"), link.WorkspaceID),
			http.StatusPermanentRedirect)
	})

	v1Router := chi.NewMux()
	v1Router.Use(getUserSession)

	v1Router.Get("/auth/logout",
		ms.getUserMiddleware(func(w http.ResponseWriter, r *http.Request, user database.User) {
			err := gothic.Logout(w, r)
			if err != nil {
				ResponseWithError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, os.Getenv("REDIRECT_URL"), http.StatusPermanentRedirect)
		}))

	v1Router.Get("/user", ms.getUserMiddleware(handleGetUser))

	// links
	v1Router.Post("/links/create/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleCreateNewWorkspaceLink)))
	v1Router.Get("/links/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceLinks)))

	// workspace
	v1Router.Get("/workspaces", ms.getUserMiddleware(ms.handleGetUserWorkspaces))
	v1Router.Get("/workspace/members/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceMembers)))
	v1Router.Get("/workspace/role/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleGetWorkspaceWithRole)))

	v1Router.Post("/workspace/create",
		ms.getUserMiddleware(ms.handleCreateWorkspace))
	v1Router.Delete("/workspace/delete/{workspace_id}",
		ms.getUserMiddleware(ms.handleDeleteWorkspace))
	v1Router.Delete("/workspace/u_r/{to_remove}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleDeleteUserFromWorksapce)))

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

	v1Router.Patch("/task/assignee/{task_id}/{workspace_id}",
		ms.getUserMiddleware(ms.checkUserWorkspaceRoleMiddleware(ms.handleUpdateTaskAssignee)))
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

	// comments
	v1Router.Post("/comment/create/{task_id}", ms.getUserMiddleware(ms.handleCreateComment))

	v1Router.Get("/comment/{task_id}", ms.getUserMiddleware(ms.handleGetTaskComments))
	v1Router.Get("/comment/p/{comment_id}", ms.getUserMiddleware(ms.handleGetChildComments))
	v1Router.Delete("/comment/delete/{comment_id}", ms.getUserMiddleware(ms.handleDeleteComment))

	router.Mount("/api/v1", v1Router)

	return router
}
