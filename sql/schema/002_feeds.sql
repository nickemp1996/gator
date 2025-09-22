-- +goose Up
CREATE TABLE feeds (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	user_id uuid NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;