-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP EXTENSION IF EXISTS pgcrypto;