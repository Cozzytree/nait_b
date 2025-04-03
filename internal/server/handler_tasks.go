package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/Cozzytree/nait/internal/model"
	"github.com/Cozzytree/nait/internal/utilfunc"
	"github.com/google/uuid"
)

const (
	DEFAULT_LIMIT = 10
)

func (ms *my_server) handleCreateTask(w http.ResponseWriter, r *http.Request, user database.User) {
	type body struct {
		WorkspaceId uuid.UUID             `json:"workspace_id" validate:"required"`
		Assignee    uuid.NullUUID         `json:"assignee"`
		Name        string                `json:"name" validate:"required,max=100,min=3"`
		Description string                `json:"description"`
		Due         time.Time             `json:"due"`
		ParentTask  uuid.NullUUID         `json:"parent_task"`
		Status      database.TaskStatus   `json:"status"`
		Priority    database.TaskPriority `json:"priority"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utilfunc.Validate(b)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(b.Description)

	err = ms.db.CreateNewTask(r.Context(), database.CreateNewTaskParams{
		WorkspaceID: b.WorkspaceId,
		Assignee:    b.Assignee,
		CreatedBy: uuid.NullUUID{
			Valid: true,
			UUID:  user.ID,
		},
		Name: b.Name,
		Due: sql.NullTime{
			Time:  b.Due,
			Valid: !b.Due.IsZero(),
		},
		Description: sql.NullString{
			String: b.Description,
			Valid:  b.Description != "",
		},
		ParentTask: b.ParentTask,
		Status: database.NullTaskStatus{
			TaskStatus: b.Status,
			Valid:      b.Status != "",
		},
		Priority: database.NullTaskPriority{
			TaskPriority: b.Priority,
			Valid:        b.Priority != "",
		},
	})
	if err != nil {
		ResponseWithError(w, "errow while creating task", http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 201,
	})
}

func (ms *my_server) handleGetWorkspaceTasks(w http.ResponseWriter, r *http.Request, user database.User) {
	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)

	w_id := r.PathValue("workspace_id")
	workspace_id, err := uuid.Parse(w_id)
	if err != nil {
		ResponseWithError(w, "invalid workspace id", http.StatusBadRequest)
		return
	}

	tasks, err := ms.db.GetWorkspaceTasks(r.Context(), database.GetWorkspaceTasksParams{
		Offset:      int32(offset),
		WorkspaceID: workspace_id,
		Limit:       int32(limit),
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Data:   model.DatabaseTasksToTasks(tasks),
		Status: 200,
	})
}

func (ms *my_server) handleGetWorkspaceCompletedTask(
	w http.ResponseWriter, r *http.Request, _ database.User) {

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)

	w_id := r.PathValue("workspace_id")
	workspace_id, err := uuid.Parse(w_id)
	if err != nil {
		ResponseWithError(w, "invalid workspace id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	fmt.Println("role", role)

	tasks, err := ms.db.GetWorkspaceCompletedTasks(r.Context(),
		database.GetWorkspaceCompletedTasksParams{
			WorkspaceID: workspace_id,
			Offset:      offset,
			Limit:       limit,
		})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Data:   model.DatabaseTasksToTasks(tasks),
		Status: 200,
	})
}

func (ms *my_server) handleUpdateTaskDesciption(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body := struct {
		Description string `json:"description"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.UpdateTaskDescription(r.Context(), database.UpdateTaskDescriptionParams{
		Description: sql.NullString{
			String: body.Description,
			Valid:  true,
		},
		ID: task_id,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleUpdateTaskPriority(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body := struct {
		Priority database.TaskPriority `json:"priority"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.UpdateTaskPriority(r.Context(), database.UpdateTaskPriorityParams{
		Priority: database.NullTaskPriority{
			TaskPriority: body.Priority,
			Valid:        true,
		},
		ID: task.ID,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleUpdateTaskName(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body := struct {
		Name string `json:"name" validate:"required,max=100,min=3"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utilfunc.Validate(body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.UpdateTaskName(r.Context(), database.UpdateTaskNameParams{
		Name: body.Name,
		ID:   task.ID,
	})

	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleUpdateTaskStatus(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body := struct {
		Status database.TaskStatus `json:"status"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ms.db.UpdateTaskStatus(r.Context(), database.UpdateTaskStatusParams{
		Status: database.NullTaskStatus{
			TaskStatus: body.Status,
			Valid:      true,
		},
		ID: task.ID,
	})

	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleUpdateTaskDue(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body := struct {
		Due time.Time `json:"due" validate:"required"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utilfunc.Validate(body)
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Due.Before(time.Now()) {
		ResponseWithError(w, "due cannot be in the past", http.StatusBadRequest)
		return
	}

	err = ms.db.UpdateTaskDue(r.Context(), database.UpdateTaskDueParams{
		Due: sql.NullTime{
			Time:  body.Due,
			Valid: true,
		},
		ID: task.ID,
	})
	if err != nil {
		ResponseWithError(w, "error while updating task", http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleGetTask(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("user_role")
	_, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if err != nil {
		ResponseWithError(w, "task not found", http.StatusUnauthorized)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseTaskToTask(task),
		Status: 200,
	})
}

func (ms *my_server) handleGetChildTasks(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)

	tasks, err := ms.db.GetChildTasks(r.Context(), database.GetChildTasksParams{
		ParentTask: uuid.NullUUID{
			UUID:  task_id,
			Valid: true,
		},
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		ResponseWithError(w, "error getting tasks", http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseTasksToTasks(tasks),
		Status: 200,
	})
}

func (ms *my_server) handleDeleteTask(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	id := r.PathValue("task_id")
	task_id, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, "invalid task id", http.StatusBadRequest)
		return
	}

	// get user role from the context
	role := r.Context().Value("user_role")
	user_role, ok := role.(database.Roles)
	if !ok {
		ResponseWithError(w, "invalid user role", http.StatusBadRequest)
		return
	}

	// getting task to check if assigned
	task, err := ms.db.GetTaskById(r.Context(), task_id)
	if err != nil {
		if err == sql.ErrNoRows {
			ResponseWithError(w, "task not found", http.StatusNotFound)
		} else {
			ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if user_role == database.RolesMember && task.Assignee.UUID != user.ID {
		ResponseWithError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = ms.db.DeleteTask(r.Context(), task.ID)
	if err != nil {
		ResponseWithError(w, "error deleting task", http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Status: 200,
	})
}

func (ms *my_server) handleGetTotalTaskCount(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	w_id := r.PathValue("workspace_id")
	workspace_id, _ := uuid.Parse(w_id)

	// utilfunc.Validate(val any)
	count, err := ms.db.GetWorkspaceTotalCountTask(r.Context(), workspace_id)
	if err != nil {
		ResponseWithError(w, "error while getting count", http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Data:   count,
		Status: 200,
	})
}

func (ms *my_server) handleGetWorkspaceTaskStatusCount(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	w_id := r.PathValue("workspace_id")
	workspace_id, _ := uuid.Parse(w_id)

	s := r.PathValue("status")

	status := database.TaskStatus(s)

	// utilfunc.Validate(val any)

	count, err := ms.db.GetWorkspaceTaskStatusCount(r.Context(), database.GetWorkspaceTaskStatusCountParams{
		Status: database.NullTaskStatus{
			TaskStatus: status,
			Valid:      true,
		},
		WorkspaceID: workspace_id,
	})

	if err != nil {
		ResponseWithError(w, "error while getting count", http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Data:   count,
		Status: 200,
	})
}

func (ms *my_server) handleGetWorkspaceTaskPriorityCount(
	w http.ResponseWriter, r *http.Request, _ database.User,
) {
	w_id := r.PathValue("workspace_id")
	workspace_id, _ := uuid.Parse(w_id)

	s := r.PathValue("priority")

	priority := database.TaskPriority(s)

	// utilfunc.Validate(val any)

	count, err := ms.db.GetWorkspaceTaskPriorityCount(r.Context(), database.GetWorkspaceTaskPriorityCountParams{
		Priority: database.NullTaskPriority{
			TaskPriority: priority,
			Valid:        true,
		},
		WorkspaceID: workspace_id,
	})

	if err != nil {
		ResponseWithError(w, "error while getting count", http.StatusInternalServerError)
		return
	}
	ResponseWithJson(w, Response{
		Data:   count,
		Status: 200,
	})
}

func (ms *my_server) handleGetWorkspaceUsersAssignedTask(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	w_id := r.PathValue("workspace_id")
	workspace_id, _ := uuid.Parse(w_id)

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)
	assignedTasks, err := ms.db.GetWorkspaceUserAssignedTasks(r.Context(), database.GetWorkspaceUserAssignedTasksParams{
		Assignee: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		WorkspaceID: workspace_id,
		Offset:      offset,
		Limit:       limit,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseTasksToTasks(assignedTasks),
		Status: 200,
	})
}

func (ms *my_server) handleGetWorkspaceUsersCreatedTask(
	w http.ResponseWriter, r *http.Request, user database.User,
) {
	w_id := r.PathValue("workspace_id")
	workspace_id, _ := uuid.Parse(w_id)

	offset, limit := utilfunc.GetOffsetAndLimitFromReq(r)
	assignedTasks, err := ms.db.GetWorkspaceUserCreatedTasks(r.Context(), database.GetWorkspaceUserCreatedTasksParams{
		CreatedBy: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		WorkspaceID: workspace_id,
		Offset:      offset,
		Limit:       limit,
	})
	if err != nil {
		ResponseWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, Response{
		Data:   model.DatabaseTasksToTasks(assignedTasks),
		Status: 200,
	})
}
