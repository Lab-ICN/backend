-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email CITEXT UNIQUE NOT NULL,
  username CITEXT UNIQUE NOT NULL,
  fullname CITEXT NOT NULL,
  is_member BOOLEAN NOT NULL,
  internship_start_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
