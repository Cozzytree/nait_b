-- name: CreateNewPage :exec
INSERT INTO pages (name, workspace_id, icon, cover_image)
VALUES ($1, $2, $3, $4);

-- name: GetWorkspacePages :many
SELECT id, name, workspace_id FROM pages
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: DeletePage :exec
DELETE FROM pages WHERE id = $1;

-- name: GetPageByID :one
SELECT * FROM pages
WHERE id = $1;
