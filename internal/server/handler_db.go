package server

import (
	"encoding/json"
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/google/uuid"
)

func (ms *my_server) handleCreateNewDB(w http.ResponseWriter, r *http.Request, _ database.User) {
	createTable := struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		WorkspaceID uuid.UUID `json:"workspace_id"`
		Icon        string    `json:"icon"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&createTable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.CreateNewDBTable(r.Context(), database.CreateNewDBTableParams{
		Name:        createTable.Name,
		Description: createTable.Description,
		WorkspaceID: createTable.WorkspaceID,
		Icon:        createTable.Icon,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}

func (ms *my_server) handleCreateDBCol(w http.ResponseWriter, r *http.Request, _ database.User) {
	id := r.PathValue("table_id")
	table_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	tableField := struct {
		Name     string            `json:"name"`
		DataType database.Datatype `json:"data_type"`
	}{}

	err = ms.db.CreateNewDBCol(r.Context(), database.CreateNewDBColParams{
		TableID:  table_id,
		Name:     tableField.Name,
		DataType: tableField.DataType,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}

func (ms *my_server) handleCreateDB_Data(w http.ResponseWriter, r *http.Request, _ database.User) {
	id := r.PathValue("table_id")
	cid := r.PathValue("col_id")

	table_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "Invalid table ID", http.StatusBadRequest)
		return
	}

	col_id, err := uuid.Parse(cid)
	if err != nil {
		ResponseWithError(w, "Invalid column ID", http.StatusBadRequest)
		return
	}

	cols, err := ms.db.GetDBCols(r.Context(), table_id)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type val struct {
		Value database.Datatype `json:"value"`
	}

	var value val
	err = json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, dc := range cols {
		if dc.ID == col_id {
			ms.db.CreateColData(r.Context(), database.CreateColDataParams{
				ColID: col_id,
				Value: string(value.Value),
			})
		} else {

			ms.db.CreateColData(r.Context(), database.CreateColDataParams{
				ColID: col_id,
				Value: "",
			})
		}
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}
