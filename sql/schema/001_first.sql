-- +goose Up
CREATE TYPE roles AS ENUM ('member', 'owner', 'admin');
CREATE TYPE task_status AS ENUM ('completed', 'in-progress', 'backlog', 'review');
CREATE TYPE task_priority AS ENUM ('high', 'low', 'medium');

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  email VARCHAR(255) UNIQUE NOT NULL UNIQUE,
  auth_id TEXT NOT NULL,
  name VARCHAR(255) NOT NULL,
  avatar TEXT DEFAULT '',
  provider TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workspaces (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID NOT NULL REFERENCES users ON DELETE RESTRICT,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workspace_members (
  user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  workspace_id UUID NOT NULL REFERENCES workspaces ON DELETE CASCADE,
  role roles NOT NULL DEFAULT 'member',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, workspace_id)
);

CREATE TABLE tasks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  workspace_id UUID NOT NULL REFERENCES workspaces ON DELETE CASCADE,
  assignee UUID REFERENCES users ON DELETE SET NULL,
  created_by UUID REFERENCES users ON DELETE SET NULL,
  name TEXT NOT NULL,
  description TEXT DEFAULT '',
  due DATE,
  parent_task UUID REFERENCES tasks ON DELETE CASCADE,
  status task_status DEFAULT 'backlog',
  priority task_priority DEFAULT 'medium',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  by UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  parent_comment UUID REFERENCES comments ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  task_id UUID NOT NULL REFERENCES tasks ON DELETE CASCADE
);

CREATE TABLE pages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  workspace_id UUID NOT NULL REFERENCES workspaces ON DELETE CASCADE,
  name TEXT NOT NULL DEFAULT 'untitled',
  icon TEXT DEFAULT '',
  cover_image TEXT DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE block (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  page_id UUID NOT NULL REFERENCES pages ON DELETE CASCADE,
  block_id TEXT NOT NULL,
  props JSONB DEFAULT '{}',
  type TEXT DEFAULT '',
  content TEXT[],
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE OR REPLACE FUNCTION check_workspace_member()
-- RETURNS TRIGGER AS $$
-- BEGIN
--   -- Check if the user is a member of the workspace
--   IF NOT EXISTS (
--     SELECT 1
--     FROM workspace_members
--     WHERE workspace_id = NEW.workspace_id
--       AND user_id = NEW.created_by
--       AND (role = 'admin' OR role = 'owner')
--   ) THEN
--     -- Raise an error if the user is not a member
--     RAISE EXCEPTION 'User % is not a member of workspace %', NEW.created_by, NEW.workspace_id;
--   END IF;

--   -- If the user is a member, allow the insert/update to proceed
--   RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER validate_workspace_member
-- BEFORE INSERT OR UPDATE ON tasks
-- FOR EACH ROW
-- EXECUTE FUNCTION check_workspace_member();

-- +goose Down

-- DROP TRIGGER validate_workspace_member ON tasks;
-- DROP FUNCTION validate_workspace_member;
DROP TABLE block;

DROP TABLE pages;

DROP TABLE tasks;

DROP TABLE comments;

DROP TABLE workspace_members;

DROP TABLE workspaces;

DROP TABLE users;

DROP TYPE task_priority;
DROP TYPE task_status;
DROP TYPE roles;
