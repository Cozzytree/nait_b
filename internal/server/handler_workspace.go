package server

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/google/uuid"
)

func (ms *my_server) handleCreateWorkspace(w http.ResponseWriter, r *http.Request, user database.User) {
	body := struct {
		Name string `json:"name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(r.Body, err)
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	created_id, err := ms.db.CreateNewWorkspace(r.Context(), database.CreateNewWorkspaceParams{
		Name:   body.Name,
		UserID: user.ID,
	})
	if err != nil {
		ResponseWithError(w, "error while creating workspace", http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   created_id,
		Status: 201,
	})
}

func (ms *my_server) handleGetUserWorkspaces(w http.ResponseWriter, r *http.Request, user database.User) {
	user_workspaces, err := ms.db.GetUserWorkspaces(r.Context(), user.ID)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.UserWorkspaceFromDatabaseRows(user_workspaces),
		Status: 200,
	})
}

func (ms *my_server) handleDeleteWorkspace(w http.ResponseWriter, r *http.Request, user database.User) {
	id := r.PathValue("workspace_id")
	workspace_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.DeleteWorkspace(r.Context(), database.DeleteWorkspaceParams{
		ID:     workspace_id,
		UserID: user.ID,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   nil,
		Status: 200,
	})
}
