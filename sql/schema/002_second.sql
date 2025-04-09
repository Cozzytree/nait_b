-- +goose Up
CREATE TABLE links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  workspace_id UUID NOT NULL REFERENCES workspaces ON DELETE CASCADE,
  valid_until TIMESTAMP NOT NULL DEFAULT (NOW () + INTERVAL '1 hour'),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  link TEXT NOT NULL,
  role roles NOT NULL DEFAULT 'member'
);

-- +goose Down
DROP TABLE links;
