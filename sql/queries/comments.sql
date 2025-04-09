-- name: CreateNewComment :exec
INSERT INTO comments (by, parent_comment, content, task_id)
VALUES ($1, $2, $3, $4);

-- name: GetTaskComments :many
SELECT c.* , jsonb_build_object(
  'user_id', u.id,
  'username', u.name,
  'avatar', u.avatar,
  'email', u.email
) AS user
FROM comments AS c
LEFT JOIN users as u ON c.by = u.id
WHERE task_id = $1 AND parent_comment IS NULL
OFFSET $2 LIMIT $3;

-- name: GetChildComments :many
SELECT c.* , jsonb_build_object(
  'user_id', u.id,
  'username', u.name,
  'avatar', u.avatar,
  'email', u.email
) AS user
FROM comments AS c
LEFT JOIN users as u ON c.by = u.id
WHERE c.parent_comment = $1
OFFSET $2 LIMIT $3;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1 AND by = $2;

-- name: UpdateComment :exec
UPDATE comments
SET content = $1,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $2;
