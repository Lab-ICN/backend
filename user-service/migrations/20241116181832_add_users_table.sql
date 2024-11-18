-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  "id" BIGSERIAL PRIMARY KEY,
  "email" TEXT UNIQUE NOT NULL,
  "username" TEXT UNIQUE NOT NULL,
  "fullname" TEXT NOT NULL,
  "is_member" BOOLEAN NOT NULL,
  "internship_start_date" DATE NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
