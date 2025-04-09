-- name: GenerateWorkspaceJoinLink :one
INSERT INTO links (
  workspace_id,
  valid_until,
  link,
  role,
  id)
VALUES($1, $2, $3, $4, $5) RETURNING *;

-- name: GetActiveWorkspaceLinks :many
SELECT * FROM links
WHERE workspace_id = $1 AND valid_until > NOW()
ORDER BY valid_until DESC;

-- name: DeleteExpiredWorkspaceLinks :exec
DELETE FROM links
WHERE valid_until < NOW();

-- name: GetAlink :one
SELECT * FROM links
WHERE id = $1;
