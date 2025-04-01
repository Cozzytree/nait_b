-- name: CreateNewTask :exec
INSERT INTO
  tasks (
    workspace_id,
    assignee,
    created_by,
    name,
    description,
    due,
    parent_task,
    status,
    priority
  )
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: GetWorkspaceTasks :many
SELECT * FROM tasks
WHERE workspace_id = $1 AND status != 'completed'
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

-- name: GetTaskById :one
SELECT * FROM tasks
WHERE id = $1;

-- name: GetChildTasks :many
SELECT * FROM tasks
WHERE parent_task = $1
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

-- nmae: GetWorkspaceCompletedTasks :many
SELECT * FROM tasks
WHERE workspace_id = $1 AND status = 'completed';

-- name: GetWorkspaceDueTasks :many
SELECT * FROM tasks
WHERE due > CURRENT_DATE
ORDER BY created_at DESC
OFFSET $1
LIMIT $2;

-- name: GetWorkspaceUserAssignedTasks :many
SELECT * FROM tasks
WHERE assignee = $1 AND workspace_id = $2
ORDER BY created_at DESC
OFFSET $3
LIMIT $4;

-- name: GetWorkspaceUserCreatedTasks :many
SELECT * FROM tasks
WHERE created_by = $1 AND workspace_id = $2
ORDER BY created_at DESC
OFFSET $3
LIMIT $4;
