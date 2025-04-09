package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/google/uuid"
)

func (ms *my_server) handleCreatePage(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name        string    `json:"name"`
		WorkspaceID uuid.UUID `json:"workspace_id"`
		Icon        string    `json:"icon"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace_id, err := uuid.Parse(body.WorkspaceID.String())
	if err != nil {
		http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
		return
	}

	err = ms.db.CreateNewPage(r.Context(), database.CreateNewPageParams{
		Name:        body.Name,
		WorkspaceID: workspace_id,
		Icon: sql.NullString{
			String: body.Icon,
			Valid:  body.Icon != "",
		},
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   "page created",
		Status: 201,
	})
}

func (ms *my_server) handleGetWorkspacePages(w http.ResponseWriter, r *http.Request, user database.User) {
	workspace_id, err := uuid.Parse(r.PathValue("workspace_id"))
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace, err := ms.db.GetWorkspaceByID(r.Context(), workspace_id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Workspace not found", http.StatusNotFound)
		} else {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	pages, err := ms.db.GetWorkspacePages(r.Context(), workspace.ID)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseWorkspacePagesToWorkspacePages(pages),
		Status: 200,
	})
}

func (ms *my_server) handleDeletePage(w http.ResponseWriter, r *http.Request, user database.User) {
	workspace_id, err := uuid.Parse(r.PathValue("workspace_id"))
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p_id := r.PathValue("page_id")
	page_id, err := uuid.Parse(p_id)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	page, err := ms.db.GetPageByID(r.Context(), page_id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Page not found", http.StatusNotFound)
		} else {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	workspace, err := ms.db.GetWorkspaceByID(r.Context(), workspace_id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Workspace not found", http.StatusNotFound)
		} else {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if page.WorkspaceID != workspace.ID {
		http.Error(w, "Unknown page", http.StatusNotFound)
		return
	}

	err = ms.db.DeletePage(r.Context(), page.ID)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 200,
	})
}
