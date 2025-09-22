-- +goose Up
CREATE TABLE feed_follows (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	user_id uuid NOT NULL,
	feed_id uuid NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
        ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;