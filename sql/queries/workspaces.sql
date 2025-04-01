-- name: CreateNewWorkspace :one
INSERT INTO workspaces (name, user_id)
VALUES ($1, $2) RETURNING id;

-- name: GetUserWorkspaces :many
SELECT * FROM workspaces
where user_id = $1;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces WHERE id = $1 AND user_id = $2;

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = $1 AND user_id = $2;

-- name: CreateNewWorkspaceMember :exec
INSERT INTO workspace_members (user_id, workspace_id, role)
VALUES ($1, $2, $3);

-- name: RemoveMemberFromWorkspace :exec
DELETE FROM workspace_members
WHERE workspace_id = $1 AND user_id = $2;

-- name: GetWorkspaceMembers :many
SELECT mem.*,
jsonb_build_object(
     'user_id', u.id,
     'username', u.name,
     'avatar', u.avatar,
     'email', u.email
 ) AS user
FROM workspace_members AS mem
LEFT JOIN users as u ON mem.user_id = u.id
WHERE workspace_id = $1
ORDER BY mem.created_at DESC;
