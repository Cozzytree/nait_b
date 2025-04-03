package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Cozzytree/nait/internal/database"
	"github.com/google/uuid"
)

type GetWorkspacePages struct {
	ID          uuid.UUID
	Name        string
	WorkspaceID uuid.UUID
}

type Task struct {
	ID          uuid.UUID             `json:"id"`
	WorkspaceID uuid.UUID             `json:"workspace_id"`
	Assignee    uuid.UUID             `json:"assignee,omitzero"`
	CreatedBy   uuid.UUID             `json:"created_by"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Due         string                `json:"due"`
	ParentTask  uuid.NullUUID         `json:"parent_task"`
	Status      database.TaskStatus   `json:"status"`
	Priority    database.TaskPriority `json:"priority"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type Page struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	Name        string
	Icon        sql.NullString
	CoverImage  sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type props struct {
	BackgroundColor string `json:"backgroundColor"`
	TextAlignment   string `json:"textAlignment"`
	TextColor       string `json:"textColor"`
}

type UserRow struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AuthID    string    `json:"auth_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Avatar    string    `json:"avatar"`
}

type Block struct {
	Children []Block `json:"children"`
	Content  []any   `json:"content"`
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Props    struct {
		BackgroundColor string `json:"backgroundColor"`
		TextAlignment   string `json:"textAlignment"`
		TextColor       string `json:"textColor"`
	} `json:"props"`
}

type BlockPacket struct {
	Block   []Block   `json:"blocks"`
	User_Id uuid.UUID `json:"user_id"`
}

type UserWorkspace struct {
	WorkspaceID uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
}

type WorkspaceMembers struct {
	UserID      uuid.UUID       `json:"id"`
	WorkspaceID uuid.UUID       `json:"workspace_id"`
	Role        database.Roles  `json:"role"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	User        json.RawMessage `json:"user"`
}

func DatabaseWorkspaceMtoWorkspaceM(m []database.GetWorkspaceMembersRow) []WorkspaceMembers {
	mems := make([]WorkspaceMembers, len(m))
	for i, gwmr := range m {
		mems[i] = WorkspaceMembers{
			UserID:      gwmr.UserID,
			WorkspaceID: gwmr.WorkspaceID,
			Role:        gwmr.Role,
			CreatedAt:   gwmr.CreatedAt,
			UpdatedAt:   gwmr.UpdatedAt,
			User:        gwmr.User,
		}
	}
	return mems
}

func UserWorkspaceFromDatabaseRows(rows []database.Workspace) []UserWorkspace {
	var userWorkspaces []UserWorkspace
	for _, row := range rows {
		userWorkspaces = append(userWorkspaces, UserWorkspace{
			WorkspaceID: row.ID,
			UserID:      row.UserID,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			Name:        row.Name,
		})
	}
	return userWorkspaces
}

func DatabaseUserRowToUser(user database.User) UserRow {
	return UserRow{
		Avatar:    user.Avatar.String,
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AuthID:    user.AuthID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func DatabasePageToPages(db_pages []database.Page) []Page {
	pages := make([]Page, len(db_pages))
	for i, page := range db_pages {
		pages[i] = Page{
			ID:          page.ID,
			WorkspaceID: page.WorkspaceID,
			Name:        page.Name,
			Icon:        page.Icon,
			CoverImage:  page.CoverImage,
			CreatedAt:   page.CreatedAt,
			UpdatedAt:   page.UpdatedAt,
		}
	}
	return pages
}

func DatabaseWorkspacePagesToWorkspacePages(pages []database.GetWorkspacePagesRow) []GetWorkspacePages {
	workspacePages := make([]GetWorkspacePages, len(pages))
	for i, page := range pages {
		workspacePages[i] = GetWorkspacePages{
			ID:          page.ID,
			WorkspaceID: page.WorkspaceID,
			Name:        page.Name,
		}
	}
	return workspacePages
}

func formatDueDate(d sql.NullTime) string {
	if d.Valid {
		return d.Time.Format("01-02-2006") // Format as MM-DD-YYYY
	}
	return "" // Return empty if invalid
}

func DatabaseTasksToTasks(db_tasks []database.Task) []Task {
	tasks := make([]Task, len(db_tasks))
	for i, task := range db_tasks {
		var assign uuid.NullUUID
		if task.Assignee.Valid {
			assign.UUID = task.Assignee.UUID
		} else {
			assign = uuid.NullUUID{}
		}

		tasks[i] = Task{
			ID:          task.ID,
			WorkspaceID: task.WorkspaceID,
			Name:        task.Name,
			Assignee:    assign.UUID,
			CreatedBy:   task.CreatedBy.UUID,
			Description: task.Description.String,
			Due:         formatDueDate(task.Due),
			ParentTask:  task.ParentTask,
			Status:      task.Status.TaskStatus,
			Priority:    task.Priority.TaskPriority,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}
	return tasks
}

func DatabaseTaskToTask(task database.Task) Task {
	var assign uuid.NullUUID
	if task.Assignee.Valid {
		assign.UUID = task.Assignee.UUID
	} else {
		assign = uuid.NullUUID{}
	}
	return Task{
		ID:          task.ID,
		WorkspaceID: task.WorkspaceID,
		Name:        task.Name,
		Assignee:    assign.UUID,
		CreatedBy:   task.CreatedBy.UUID,
		Description: task.Description.String,
		Due:         formatDueDate(task.Due),
		ParentTask:  task.ParentTask,
		Status:      task.Status.TaskStatus,
		Priority:    task.Priority.TaskPriority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
