-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  published_at TIMESTAMP,
  feed_id UUID NOT NULL,

  FOREIGN KEY(feed_id) REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
