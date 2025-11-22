-- +goose Up
CREATE TABLE feeds(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_fetched_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_feed FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
