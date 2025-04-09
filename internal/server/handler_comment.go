package server

import (
	"encoding/json"
	"net/http"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/Cozzytree/nait/internal/utilfunc"
	"github.com/google/uuid"
)

func (ms *my_server) handleCreateComment(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	type body struct {
		Content       string        `json:"content" validate:"required,max=1000"`
		ParentComment uuid.NullUUID `json:"parent_comment"`
	}

	var b body
	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utilfunc.Validate(b)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.CreateNewComment(r.Context(), database.CreateNewCommentParams{
		By:            user.ID,
		ParentComment: b.ParentComment,
		Content:       b.Content,
		TaskID:        task_id,
	})
	if err != nil {
		ResponseWithError(w, "error while creating task", http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}

func (ms *my_server) handleGetTaskComments(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)

	comments, err := ms.db.GetTaskComments(r.Context(), database.GetTaskCommentsParams{
		TaskID: task_id,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 200,
		Data:   model.DatabaseCommentsToComments(comments),
	})
}

func (ms *my_server) handleGetChildComments(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	id := r.PathValue("comment_id")
	comment_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)

	comments, err := ms.db.GetChildComments(r.Context(), database.GetChildCommentsParams{
		ParentComment: uuid.NullUUID{
			UUID:  comment_id,
			Valid: true,
		},
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 200,
		Data:   model.DatabaseChildCommentsToComments(comments),
	})
}

func (ms *my_server) handleDeleteComment(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("comment_id")
	comment_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	err = ms.db.DeleteComment(r.Context(), database.DeleteCommentParams{
		ID: comment_id,
		By: user.ID,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 200,
	})
}
