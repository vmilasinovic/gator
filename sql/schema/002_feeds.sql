-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE feeds;