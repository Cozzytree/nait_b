package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/google/uuid"
)

func (ms *my_server) handleCreateNewWorkspaceLink(w http.ResponseWriter, r *http.Request, _ database.User) {
	workspace_id, _ := uuid.Parse(r.PathValue("workspace_id"))

	body := struct {
		ValidUntil time.Time      `json:"valid_until"`
		Role       database.Roles `json:"role"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newId := uuid.New()

	_, err = ms.db.GenerateWorkspaceJoinLink(r.Context(), database.GenerateWorkspaceJoinLinkParams{
		ID:          newId,
		Role:        body.Role,
		WorkspaceID: workspace_id,
		ValidUntil:  body.ValidUntil,
		Link:        fmt.Sprintf("%s/%s/%v", os.Getenv("APP_URL"), "invite", newId),
	})

	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}

func (ms *my_server) handleGetWorkspaceLinks(w http.ResponseWriter, r *http.Request, _ database.User) {
	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	if user_role == database.RolesMember {
		ResponseWithError(w, "not authorized", http.StatusForbidden)
		return
	}

	workspace_id, _ := uuid.Parse(r.PathValue("workspace_id"))

	links, err := ms.db.GetActiveWorkspaceLinks(r.Context(), workspace_id)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseLinksToLinks(links),
		Status: 200,
	})
}
