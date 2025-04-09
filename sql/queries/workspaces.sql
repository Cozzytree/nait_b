-- name: CreateNewWorkspace :one
INSERT INTO workspaces (name, user_id)
VALUES ($1, $2) RETURNING id;

-- name: GetUserWorkspaces :many
SELECT mem.* , jsonb_build_object(
      'workspace_id', ws.id,
      'name', ws.name,
      'created_at', ws.created_at,
      'updated_at', ws.updated_at
  ) AS workspace FROM workspace_members AS mem
LEFT JOIN workspaces AS ws ON mem.workspace_id = ws.id
WHERE mem.user_id = $1
ORDER BY mem.created_at;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces WHERE id = $1 AND user_id = $2;

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = $1;

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

-- name: GetWorkspaceUserRole :one
SELECT * FROM workspace_members
WHERE user_id = $1 AND workspace_id = $2;

-- name: DeleteUserFromWorksapce :exec
DELETE FROM workspace_members
WHERE user_id = $1 AND workspace_id = $2;


-- CREATE OR REPLACE FUNCTION set_null_on_user_removal()
-- RETURNS TRIGGER AS $$
-- BEGIN
--   -- Set assignee and created_by to NULL in tasks when the user is removed from the workspace
--   UPDATE tasks
--   SET assignee = NULL,
--       created_by = NULL
--   WHERE assignee = OLD.user_id OR created_by = OLD.user_id;

--   RETURN OLD;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_set_null_on_user_removal
-- AFTER DELETE ON workspace_members
-- FOR EACH ROW
-- EXECUTE FUNCTION set_null_on_user_removal();
