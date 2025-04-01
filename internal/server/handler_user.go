package server

import (
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
)

func handleGetUser(w http.ResponseWriter, _ *http.Request, user database.User) {
	ResponseWithJson(w, Response{
		Data:   model.DatabaseUserRowToUser(user),
		Status: 200,
	})
}

// func handleCreateUser(w http.ResponseWriter)
