-- nmae: GetCountOfWorkspaceCompletdTask :one
SELECT COUNT(id) FROM tasks
WHERE workspace_id = $1 AND status = 'completed';

-- name: GetWorkspaceUserAssignedTaskCount :one
SELECT COUNT(id) FROM tasks
WHERE assignee IS NOT NULL AND workspace_id = $1 AND assignee = $2;

-- name: GetWorkspaceUserCreatedTaskCount :one
SELECT COUNT(id) FROM tasks
WHERE created_by IS NOT NULL AND workspace_id = $1 AND created_by = $2;

-- name: GetWorkspaceTotalCountTask :one
SELECT COUNT(id) FROM tasks
WHERE workspace_id = $1;

-- name: GetWorkspaceTaskStatusCount :one
SELECT COUNT(id) FROM tasks
WHERE status = $1 AND workspace_id = $2;

-- name: GetWorkspaceTaskPriorityCount :one
SELECT COUNT(id) FROM tasks
WHERE priority = $1 AND workspace_id = $2;

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
WHERE workspace_id = $1
AND status != 'completed'
AND parent_task IS NULL
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

-- name: GetTaskById :one
SELECT t.*, jsonb_build_object(
           'username', au.name,
           'email', au.email,
           'avatar', au.avatar) AS assigned,
           jsonb_build_object(
           'username', cu.name,
           'email', cu.email,
           'avatar', cu.avatar) AS created
FROM tasks AS t
LEFT JOIN users AS au ON t.assignee = au.id
LEFT JOIN users AS cu ON t.created_by = cu.id
WHERE t.id = $1;


-- name: GetChildTasks :many
SELECT * FROM tasks
WHERE parent_task = $1
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

-- name: GetWorkspaceCompletedTasks :many
SELECT * FROM tasks
WHERE workspace_id = $1 AND status = 'completed'
ORDER BY created_at DESC
OFFSET $2
LIMIT $3;

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

-- name: UpdateTaskDue :exec
UPDATE tasks
SET due = $1,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $2 AND status != 'completed';

-- name: UpdateTaskDescription :exec
UPDATE tasks
SET description = $1,
updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateTaskPriority :exec
UPDATE tasks
SET priority = $1,
updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateTaskName :exec
UPDATE tasks
SET name = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateTaskAssignee :exec
UPDATE tasks
SET assignee = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateTaskCreated :exec
UPDATE tasks
SET created_by = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;
