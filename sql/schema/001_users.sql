-- +goose Up 
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  update_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
